//Package stata writes data into a Stata 113 format (readable by any Stata version higher than 7)
//Source for format info https://www.stata.com/help.cgi?dta_113
//The package does not do much validation. It is up to the user to ensure that the supplied data
//meets the format specification!
package stata

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"unsafe"
)

const (
	stataVarSize   = 33
	stataFmtSize   = 12
	stataLabelSize = 81

	STATA_BYTE_NA     = 127
	STATA_SHORTINT_NA = 32767
	STATA_INT_NA      = 2147483647
)

var (
	STATA_FLOAT_NA  = math.Pow(2.0, 127)
	STATA_DOUBLE_NA = math.Pow(2.0, 1023)

	littleEndian = binary.LittleEndian //only write LittleEndian
)

type (
	stataVarName [stataVarSize]byte //used for both variable and label names
	stataFmtName [stataFmtSize]byte //used for format names
	stataLabel   [stataLabelSize]byte
)

//Supported Stata var types aliased to Go types for maximum convertability
type (
	Byte   = int8
	Int    = int16
	Long   = int32
	Float  = float32
	Double = float64
)

//Supported Stata variable types
const (
	StataByteId   = 251 // 0xfb
	StataIntId    = 252 // 0xfc
	StataLongId   = 253 // 0xfd
	StataFloatId  = 254 // 0xfe
	StataDoubleId = 255 // 0xff
)

/*
         type          code
                --------------------
                str1        1 = 0x01
                str2        2 = 0x02
                ...
                str244    244 = 0xf4
                byte      251 = 0xfb  (sic)
                int       252 = 0xfc
                long      253 = 0xfd
                float     254 = 0xfe
                double    255 = 0xff
				--------------------
*/

//Field holds information a Stata variable
type Field struct {
	Name      string
	FieldType byte
	Label     string
	Format    string
	data      interface{}
}

//field name must be exported for package Binary to see them
type header struct {
	//	Contents            Length    Format    Comments
	Version   byte       //     1    byte      contains 113 = 0x71
	ByteOrder byte       //     1    byte      0x01 -> HILO, 0x02 -> LOHI
	FileType  byte       //     1    byte      0x01
	UnUsed    byte       //     1    byte      0x01
	NoVar     int16      //		2    int       (number of vars) encoded per byteorder
	NoObs     int32      // 		4    int       (number of obs)  encoded per byteorder
	DataLabel stataLabel // 		81    char      dataset label, \0 terminated
	TimeStamp [18]byte   // 		18    char      date/time saved, \0 terminated
	//	Total                  109
}

func NewHeader() *header {
	fh := header{
		Version:   113, //113 is used in Stata versions 8-9
		ByteOrder: 2,   //LOHI
		FileType:  1,   //always 1
		UnUsed:    0,
	}
	//FIXME: leave empty for production; comment the line below
	copy(fh.DataLabel[:], "Written by VDEC Stata File Creator")
	// fh.timeStamp[0] = 0
	return &fh
}

//File Stata file info
type File struct {
	*header
	fields     []*Field
	recordSize int
	//FIXME: remove from the struct and just declare when needed?
	//	Contents            	Length    	  Format       Comments
	typList  []byte         //         nvar    byte array
	varList  []stataVarName //      33*nvar    char array
	srtList  []byte         //    2*(nvar+1)   int array    encoded per byteorder
	fmtList  []stataFmtName //      12*nvar    char array
	lblList  []stataVarName //       33*nvar    char array
	vlblList []stataLabel
}

//NewFile returns a pointer to an initialized File.
func NewFile() *File {
	sf := File{
		header: NewHeader(),
	}
	return &sf
}

//AddField adds a field to be written out to a Stata file
//It does not verify similarly-named field does not exist
//It does not verify field names and labels meet Stata requirements
//It does not verify that slice lengths are identical
func (sf *File) AddField(name, label string, slice interface{}) *Field {
	var (
		typ      byte
		sliceLen int
		format   = "%9.0g"
	)

	switch data := slice.(type) {
	case []Byte:
		typ = StataByteId
		sf.recordSize++ //one byte
		sliceLen = len(data)
	case []Int:
		typ = StataIntId
		sf.recordSize += 2
		sliceLen = len(data)
	case []Long:
		typ = StataLongId
		sf.recordSize += 4
		sliceLen = len(data)
	case []Float:
		typ = StataFloatId
		sf.recordSize += 4
		sliceLen = len(data)
	case []Double:
		typ = StataDoubleId
		sf.recordSize += 8
		sliceLen = len(data)
	default:
		panic("unsupported data type in field " + name) //must be a programmer error, so panic
		//return nil, fmt.Errorf("unsupported data type in field %s", name)
	}
	fld := &Field{
		Name:      name,
		FieldType: typ,
		Label:     label,
		Format:    format,
		data:      slice,
	}
	sf.fields = append(sf.fields, fld)
	sf.NoVar = int16(len(sf.fields))
	if sliceLen > int(sf.NoObs) {
		sf.NoObs = int32(sliceLen)
	}
	return fld
}

//WriteTo writes the data to an io.Writer.
//warning: the number of written byte is not used, always zero
func (sf *File) WriteTo(w io.Writer) (int64, error) {
	if err := sf.writeHeader(w); err != nil {
		return 0, err
	}
	if err := sf.writeDescriptors(w); err != nil {
		return 0, err
	}
	return 0, sf.writeData(w)
}

func (sf *File) writeHeader(w io.Writer) error {
	sf.NoVar = int16(len(sf.fields))
	return binary.Write(w, littleEndian, *sf.header)
}

func (sf *File) writeDescriptors(w io.Writer) error {
	sf.typList = make([]byte, sf.NoVar)
	sf.varList = make([]stataVarName, sf.NoVar)
	sf.srtList = make([]byte, 2*(sf.NoVar+1))
	sf.fmtList = make([]stataFmtName, sf.NoVar)
	sf.lblList = make([]stataVarName, sf.NoVar)
	sf.vlblList = make([]stataLabel, sf.NoVar)
	for i, f := range sf.fields {
		copy(sf.varList[i][:], f.Name) //only copy up to the size of stataVarName and pad with zeros
		sf.typList[i] = f.FieldType
		copy(sf.fmtList[i][:], f.Format)
		copy(sf.vlblList[i][:], f.Label)
	}

	if err := binary.Write(w, littleEndian, sf.typList); err != nil {
		return err
	}
	if err := binary.Write(w, littleEndian, sf.varList); err != nil {
		return err
	}
	//write an empty sort list
	if err := binary.Write(w, littleEndian, sf.srtList); err != nil {
		return err
	}
	//write var format, for now just generic numberic "%9.0g"
	if err := binary.Write(w, littleEndian, sf.fmtList); err != nil {
		return err
	}
	//write empty value lables
	if err := binary.Write(w, littleEndian, sf.lblList); err != nil {
		return err
	}
	if err := binary.Write(w, littleEndian, sf.vlblList); err != nil {
		return err
	}
	// write an empty expansion field (5 bytes of zeros)
	return binary.Write(w, littleEndian, [5]byte{0, 0, 0, 0, 0})
}

//writeData loops over the field vectors and write their binary representation to an io.Writer
func (sf *File) writeData(w io.Writer) error {
	if sf.NoObs == 0 {
		return nil
	}
	if len(sf.fields) == 0 {
		return fmt.Errorf("No fields")
	}
	bs := make([]byte, sf.recordSize)
	for i := int32(0); i < sf.NoObs; i++ {
		offset := 0
		for _, f := range sf.fields {
			switch f.FieldType {
			case StataByteId:
				v := f.data.([]Byte)[i]
				bs[offset] = byte(v)
				offset++
			case StataIntId:
				v := f.data.([]Int)[i]
				bs[offset] = byte(v)
				offset++ //incrementing the offset instead of using bs[offset+1] to avoid doing the addition twice
				bs[offset] = byte(v >> 8)
				offset++
			case StataLongId:
				base := *(*[4]byte)(unsafe.Pointer(&f.data.([]Long)[i]))
				copy(bs[offset:], base[:])
				offset += 4
			case StataFloatId:
				base := *(*[4]byte)(unsafe.Pointer(&f.data.([]Float)[i]))
				copy(bs[offset:], base[:])
				offset += 4
			case StataDoubleId:
				base := *(*[8]byte)(unsafe.Pointer(&f.data.([]Double)[i]))
				copy(bs[offset:], base[:])
				offset += 8
			default:
				return fmt.Errorf("Field type [%d] not supported in field %s", f.FieldType, f.Name)
			}
		}
		if _, err := w.Write(bs); err != nil {
			return err
		}
	}
	return nil
}

//FIXME: do not overwrite an existing file
//WriteFile
func (sf *File) WriteFile(fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	bw := bufio.NewWriterSize(f, 64*1012) //use 64kb buffer
	if _, err = sf.WriteTo(bw); err != nil {
		f.Close()
		return err
	}
	bw.Flush()
	err = f.Close()
	return err
}
