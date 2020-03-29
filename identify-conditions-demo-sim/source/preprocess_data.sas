%let LABEL_CONDITION_ID = 'Condition identification tag';
%let LABEL_DATA_SOURCE = 'Source database tag';

data processed_hosp;
		infile &hosp_file delimiter = ',' MISSOVER DSD firstobs=2;
			informat subject_id $32.;
			informat service_date YYMMDD10.;
			informat discharge_date YYMMDD10.;
			informat code $8.;
			format subject_id $32.;
			format service_date date11.;
			format discharge_date date11.;
			format code $8.;
		input subject_id $
			  service_date
			  discharge_date
			  code $;
		;
run;

data processed_phys;
		infile &phys_file delimiter = ',' MISSOVER DSD firstobs=2;
			informat subject_id $32.;
			informat service_date YYMMDD10.;
			informat code $8.;
			format subject_id $32.;
			format service_date date11.;
			format code $8.;
		input subject_id
			  service_date
			  code $;
		;
run;

data study_pop;
		infile &pop_file delimiter = ',' MISSOVER DSD firstobs=2;
			informat subject_id $32.;
			informat gender $1.;
			informat birthdate YYMMDD10.;
			informat age 8.;
			informat coverage_start YYMMDD10.;
			informat coverage_end YYMMDD10.;
			format subject_id $32.;
			format gender $1.;
			format birthdate date11.;
			format age 8.;
			format coverage_start date11.;
			format coverage_end date11.;
		input subject_id $
			  gender $
			  birthdate
			  age
			  coverage_start
			  coverage_end;
		;
run;
/*
data processed_dpin;
		infile &dpin_file delimiter = ',' MISSOVER DSD firstobs=2;
			informat subject_id $32.;
			informat service_date date11.;
			informat coding_system $4.;
			informat code $8.;
			format subject_id 32.;
			format service_date date11.;
			format coding_system $4.;
			format code $8.;
		input subject_id $
			  service_date
			  coding_system $
			  code $;
		;
run;
*/

data processed_dpin;
		infile &dpin_file delimiter = ',' MISSOVER DSD firstobs=2;
			informat subject_id $32.;
			informat service_date YYMMDD10.;
			informat code $8.;
			format subject_id 32.;
			format service_date date11.;
			format code $8.;
		input subject_id $
			  service_date
			  code $;
		;
run;

data processed_dpin;
	set processed_dpin;
	coding_system = "DIN";
run;

data diab_conditions;
	infile &cond_file delimiter = ',' MISSOVER DSD firstobs=2;
		informat condition_id $32.;
		informat condition_name $256.;
		format condition_id $32.;
		format condition_name $256.;
	input condition_id $
		  condition_name $
	;

	label condition_id = &LABEL_CONDITION_ID;
	label condition_name = 'Name of condition';
run;

* import chronic condition criteria;
data diab_criteria;
	infile &crit_file delimiter = ',' missover dsd firstobs=2;
		informat condition_id $32. ;
		informat data_source $4. ;
		informat min_num_cases best12.;
		informat time_between_cases 12.;
		informat time_period 12.;
		informat time_period_index_date $64. ;
		format condition_id $32. ;
		format data_source $4. ;
		format min_num_cases best12.;
		format time_between_cases best12.;
		format time_period best12.;
		format time_period_index_date $64. ;
	input condition_id $
		  data_source $
		  min_num_cases
		  time_between_cases
		  time_period
		  time_period_index_date $
	;

	label condition_id = 'Condition identification tag';
	label data_source = 'Source database tag';
	label min_num_cases = 'Minimum number of matching codes';
	label time_period_index_date = 'Look for condition before or after index date';
run;

* import chronic condition patterns;
data diab_patterns;
	infile &patt_file delimiter = ',' missover dsd firstobs=2;
		informat condition_id $32. ;
		informat data_source $4. ;
		informat coding_system $10. ;
		informat code_matching_pattern $10. ;
		informat description $256. ;
		format condition_id $32. ;
		format data_source $4. ;
		format coding_system $10. ;
		format code_matching_pattern $10. ;
		format description $256. ;
	input condition_id $
		  data_source $
		  coding_system $
		  code_matching_pattern $
		  description $
	;

	label condition_id = &LABEL_CONDITION_ID;
	label data_source = &LABEL_DATA_SOURCE;
	label coding_system = 'Coding system used for code';
	label code_matching_pattern = 'Match pattern (with wildcards)';
	label description = 'Description';
run;

data rxref;
		infile &rxref_file delimiter = ',' MISSOVER DSD firstobs=2;
			informat atc $7.;
			informat din $8.;
			format atc $7.;
			format din $8.;
		input atc $
			  din $
		;
run;

data processed_hosp;
	set processed_hosp;
	code = compress(code, '.');
	coding_system = "ICD-10";
run;

data processed_phys;
	set processed_phys;
	coding_system = "ICD-9";
run;

data study_pop;
	set study_pop;
	study_start_date = coverage_start;
	study_end_date   = coverage_end;
	
	format study_start_date
		   study_end_date
		   date11.;
run;