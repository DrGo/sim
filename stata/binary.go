package stata

//
// func binWrite(w io.Writer, data interface{}) error {
// 	// Fast path for basic types and slices.
// 	if n := intDataSize(data); n != 0 {
// 		var b [8]byte
// 		var bs []byte
// 		if n > len(b) {
// 			bs = make([]byte, n)
// 		} else {
// 			bs = b[:n]
// 		}
// 		switch v := data.(type) {
// 		case *bool:
// 			if *v {
// 				b[0] = 1
// 			} else {
// 				b[0] = 0
// 			}
// 		case bool:
// 			if v {
// 				b[0] = 1
// 			} else {
// 				b[0] = 0
// 			}
// 		case []bool:
// 			for i, x := range v {
// 				if x {
// 					bs[i] = 1
// 				} else {
// 					bs[i] = 0
// 				}
// 			}
// 		case *int8:
// 			b[0] = byte(*v)
// 		case int8:
// 			b[0] = byte(v)
// 		case []int8:
// 			for i, x := range v {
// 				bs[i] = byte(x)
// 			}
// 		case *uint8:
// 			b[0] = *v
// 		case uint8:
// 			b[0] = v
// 		case []uint8:
// 			bs = v
// 		case *int16:
// 			order.PutUint16(bs, uint16(*v))
// 		case int16:
// 			order.PutUint16(bs, uint16(v))
// 		case []int16:
// 			for i, x := range v {
// 				order.PutUint16(bs[2*i:], uint16(x))
// 			}
// 		case *uint16:
// 			order.PutUint16(bs, *v)
// 		case uint16:
// 			order.PutUint16(bs, v)
// 		case []uint16:
// 			for i, x := range v {
// 				order.PutUint16(bs[2*i:], x)
// 			}
// 		case *int32:
// 			order.PutUint32(bs, uint32(*v))
// 		case int32:
// 			order.PutUint32(bs, uint32(v))
// 		case []int32:
// 			for i, x := range v {
// 				order.PutUint32(bs[4*i:], uint32(x))
// 			}
// 		case *uint32:
// 			order.PutUint32(bs, *v)
// 		case uint32:
// 			order.PutUint32(bs, v)
// 		case []uint32:
// 			for i, x := range v {
// 				order.PutUint32(bs[4*i:], x)
// 			}
// 		case *int64:
// 			order.PutUint64(bs, uint64(*v))
// 		case int64:
// 			order.PutUint64(bs, uint64(v))
// 		case []int64:
// 			for i, x := range v {
// 				order.PutUint64(bs[8*i:], uint64(x))
// 			}
// 		case *uint64:
// 			order.PutUint64(bs, *v)
// 		case uint64:
// 			order.PutUint64(bs, v)
// 		case []uint64:
// 			for i, x := range v {
// 				order.PutUint64(bs[8*i:], x)
// 			}
// 		}
// 		_, err := w.Write(bs)
// 		return err
// 	}
