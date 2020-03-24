/*
Macro: 
	identify_conditions
	
-Since a single DIN code can map to more than one ATC code, it is important to ensure there is no overlap between Rx drug code patterns.

Params:
	temp_lib_: the library (work, vdec) where temporary datasets should be saved. If the vdec library is used then the datasets must be prefixed by the project short name in order to 
		       avoid concurrency issues and to avoid violating privacy regulations. Currently, temp_lib_ is only used for datasets that have the potential to be very large and problematic.
			   All other temporary datasets are still saved to the work library.
	debug_mode_: should be set to 0 if working with very large datasets. This parameter will ensure theintermediate datasets
								 are deleted when no longer needed in order to minimize the chances of running out of space in the temporary library. It will also avoid doing unnecessary and time-consuming
								 operations such as sorting that are useful when debugging.
	condition_criteria_ds_: the time_between_cases and time_period variables must be defined in units of days (positive integers). These two time variables should be set to missing
							if the min_num_cases is equal to 1 (will always be true regardless of the size of the time period).

Notes:
	-the time_period variable, defined in condition_criteria_ds_, refers to in what time period (number of days) the minimum number of cases must occur within. (ex. minimum of 4 cases within a 2 year period)
		-the time_period variable does not refer to the follow-up time prior to index date that the search will be performed in. If you want to look in the 5 year period prior to the index date
		 then this needs to be factored into the study_start_date and study_end_date variables that must be present in the provided population datsets.
		-the data preparation macro (%prepare_data_for_condition_ident) has a prepare_pop_ds_params parameter where study_start_date_expr_ and study_end_date_expr_ must be specified to factor in 
		 the time prior to index date constraint.
	-the time_between_cases variable, defined in condition_criteria_ds_, refers to the number of days that must separate instances of a disease definition being met. 
		For example, {time_period=730, time_between_cases=30, min_num_cases=5} means in order to be considered a case, a subject must have matched one of the condition codes at least 5 times in a 2 year period,
		where each instance must be at least 1 month apart.
	-this macro currently assumes a separate rule in condition_criteria_ds_ will be provided for both tariff codes and physician diagnostic codes. This macro does not handle combining these two coding systems 
	 under the same datasource (PHYS) and using a single min_num_cases against them. The TARF datasource must be used for all MH-TARIFF codes and the PHYS datasource for all ICD-9 codes in the claims data.
	 If you combine ICD-9 and MH-TARIFF codes under the PHYS datasource rule, the macro will not work correctly.

Usage:
	Normal usage: Will give you the first occurence of a condition
	
	Extended usage: It is possible to get occurences for multiple visits or complex occurences of a condition by using intermediate tables.
		This example will show how to do this for hospital data, where all occurences for each person that have either condition
		cuti_a OR (cuti_b1 AND cuti_b2) at the same time. The data is obtained from intermediate macro table Num_cases_per_subj_hosp_visit
		proc sql noprint;
	
			create table cUTI_A_visits as
			select &subject_id_var, admdt, sepdt, 'cuti' as condition_id
			from Num_cases_per_subj_hosp_visit
			where condition_id eq 'cuti_a';

			create table cUTI_B_visits as
			select &subject_id_var, admdt, sepdt, 'cuti' as condition_id
			from Num_cases_per_subj_hosp_visit
			where condition_id eq 'cuti_b1'
				or  condition_id eq 'cuti_b2'
			group by &subject_id_var, admdt, sepdt
			having count(*) > 1;

			create table cUTI_visits as 
			select * from cUTI_A_visits
			union
			select * from cUTI_B_visits;
			
		quit;

WARNING:
	-Thorough testing must be done if using DIN and ATC codes together (from DPIN dataset) or if using tariff and diagnostic codes together (from PHYS dataset). Testing these
	 situations has not yet been done.
	-Time period functionality has not been thoroughly tested (see some important test cases below).

TODO: 
	-test new time period functionality when multiple conditions are found for the same individual
	-test new time period functionality when both tariff codes and physician diagnostic codes are used with the PHYS data set
	-decide how the physician tariff codes and physician diagnostic codes should work together (especially concerning time period constraints)
	-implement functionality to keep intermediate datasets when desired, delete temp datasets otherwise;
	-validate input (ensure required datasets, fields have been provided);
*/
%macro identify_conditions(
	pop_ds_ =,
	hosp_ds_ =,
	phys_ds_ =,
	dpin_ds_ =,
	conditions_ds_ =,
	condition_matching_patterns_ds_ =,
	condition_criteria_ds_ =,
	rxdrug_reference_ds_ =,
	PHIN_var_ =,
	subject_id_var_ =,
	output_widefmt_binary_ds_ =.,
	output_widefmt_date_ds_ =.,
	output_longfmt_ds_ =.,
	temp_lib_ = work,
	debug_mode_ = 0
);

	%put %STR( -> Start of %upcase( &sysmacroname ) macro );

	%let inputs_are_valid = ;
	%validate_input_data(inputs_are_valid);

	%if &inputs_are_valid %then %do;

		%let hosp_output_ds_ = work.hosp_cases;
		%let phys_output_ds_ = work.phys_cases;
		%let dpin_output_ds_ = work.dpin_cases;

		*If the result datasets already exist (from a previous call to this macro), then they must be dropped to 
		 ensure previous results are not carried forward;
		proc datasets noprint;
			delete hosp_cases;
			delete phys_cases;
			delete dpin_cases;
		quit;

		*Determine which datasets will be used for condition identification;
		proc sql noprint;

			select (count(*) > 0)
			into :match_using_hosp_ds
			from &condition_matching_patterns_ds_
			where data_source = 'HOSP';
			
			select (count(*) > 0)
			into :match_using_phys_ds
			from &condition_matching_patterns_ds_
			where data_source = 'PHYS' or data_source = 'TARF';
			
			select (count(*) > 0)
			into :match_using_dpin_ds
			from &condition_matching_patterns_ds_
			where data_source = 'DPIN';

		quit;

		%if %eval(&match_using_hosp_ds = 1) %then %do;
		
			%identify_cases_in_hosp_ds(
				pop_ds_ = &pop_ds_,
				hosp_ds_ = &hosp_ds_,
				condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
				condition_criteria_ds_ = &condition_criteria_ds_,
				PHIN_var_ = &PHIN_var_,
				subject_id_var_ = &subject_id_var_,
				output_ds_ = &hosp_output_ds_
			);
		
		%end;

		%if %eval(&match_using_phys_ds = 1) %then %do;
		
			%identify_cases_in_phys_ds(
				pop_ds_ = &pop_ds_,
				phys_ds_ = &phys_ds_,
				condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
				condition_criteria_ds_ = &condition_criteria_ds_,
				PHIN_var_ = &PHIN_var_,
				subject_id_var_ = &subject_id_var_,
				output_ds_ = &phys_output_ds_
			);

		%end;

		%if %eval(&match_using_dpin_ds = 1) %then %do;
		
			%identify_cases_in_dpin_ds(
				pop_ds_ = &pop_ds_,
				dpin_ds_ = &dpin_ds_,
				condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
				condition_criteria_ds_ = &condition_criteria_ds_,
				rxdrug_reference_ds_ = &rxdrug_reference_ds_,
				PHIN_var_ = &PHIN_var_,
				subject_id_var_ = &subject_id_var_,
				output_ds_ = &dpin_output_ds_
			);

		%end;

		%create_final_ds(
			pop_ds_ = &pop_ds_,
			hosp_condition_cases_ = &hosp_output_ds_,
			phys_condition_cases_ = &phys_output_ds_,
			dpin_condition_cases_ = &dpin_output_ds_,
			conditions_ds_ = &conditions_ds_,
			subject_id_var_ = &subject_id_var_,
			output_widefmt_binary_ds_ = &output_widefmt_binary_ds_,
			output_widefmt_date_ds_ = &output_widefmt_date_ds_,
			output_longfmt_ds_ = &output_longfmt_ds_
		);

	%end;
	
	%put %STR( -> End of %upcase( &sysmacroname ) macro );

%mend;

**@TODO validate all input data (caller provided the required datasets and they contain the correct fields, etc.);
%macro validate_input_data(return_var);

	%put %STR( -> Start of %upcase( &sysmacroname ) macro );

	%let &return_var = 1;
	
	%put %STR( -> End of %upcase( &sysmacroname ) macro );

%mend;

%macro identify_cases_in_hosp_ds(
	pop_ds_ =,
	hosp_ds_ =,
	condition_matching_patterns_ds_ =,
	condition_criteria_ds_ =,
	PHIN_var_ =,
	subject_id_var_ =,
	output_ds_ =
);

	%put %STR( -> Start of %upcase( &sysmacroname ) macro );

	%let DATA_SOURCE_CODE = HOSP;
	
	*split the dataset name into two tokens (the library name and the table name);
	%let input_lib_name = %scan(&hosp_ds_, 1, '.');
	%let input_tbl_name = %scan(&hosp_ds_, 2, '.');

	%select_obs_by_study_date(
		input_ds_ = &hosp_ds_,
		output_ds_ = _hosp_1_ds,
		pop_ds_ = &pop_ds_,
		PHIN_var_ = &PHIN_var_,
		subject_id_var_ = &subject_id_var_
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=&input_lib_name nolist;
			delete &input_tbl_name;
		run;
	%end;

	%select_obs_with_conditions(
		input_ds_ = _hosp_1_ds,
		output_ds_ = _hosp_2_ds,
		condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
		subject_id_var_ = &subject_id_var_,
		data_source_codes_ = ("&DATA_SOURCE_CODE")
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=work nolist;
			delete _hosp_1_ds;
		run;
	%end;

	%apply_condition_criteria(
		input_ds_ = _hosp_2_ds,
		output_ds_ = &output_ds_,
		data_source_code_ = &DATA_SOURCE_CODE,
		condition_criteria_ds_ = &condition_criteria_ds_,
		subject_id_var_ = &subject_id_var_
	);

	%put %STR( -> End of %upcase( &sysmacroname ) macro );

%mend;

%macro identify_cases_in_phys_ds(
	pop_ds_ =,
	phys_ds_ =,
	condition_matching_patterns_ds_ =,
	condition_criteria_ds_ =,
	PHIN_var_ =,
	subject_id_var_ =,
	output_ds_ =
);

	%put %STR( -> Start of %upcase( &sysmacroname ) macro );
	
	%let PHYSICIAN_DATA_SOURCE_CODE = PHYS;
	%let TARIFF_DATA_SOURCE_CODE = TARF;
	
	*split the dataset name into two tokens (the library name and the table name);
	%let input_lib_name = %scan(&phys_ds_, 1, '.');
	%let input_tbl_name = %scan(&phys_ds_, 2, '.');
	
	%select_obs_by_study_date(
		input_ds_ = &phys_ds_,
		output_ds_ = _phys_1_ds,
		pop_ds_ = &pop_ds_,
		PHIN_var_ = &PHIN_var_,
		subject_id_var_ = &subject_id_var_
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=&input_lib_name nolist;
			delete &input_tbl_name;
		run;
	%end;

	%select_obs_with_conditions(
		input_ds_ = _phys_1_ds,
		output_ds_ = _phys_2_ds,
		condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
		subject_id_var_ = &subject_id_var_,
		data_source_codes_ = ("&PHYSICIAN_DATA_SOURCE_CODE", "&TARIFF_DATA_SOURCE_CODE")
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=work nolist;
			delete _phys_1_ds;
		run;
	%end;

	%apply_condition_criteria(
		input_ds_ = _phys_2_ds,
		output_ds_ = &output_ds_,
		data_source_code_ = &PHYSICIAN_DATA_SOURCE_CODE,
		condition_criteria_ds_ = &condition_criteria_ds_,
		subject_id_var_ = &subject_id_var_
	);

	%put %STR( -> End of %upcase( &sysmacroname ) macro );

%mend;

%macro identify_cases_in_dpin_ds(
	pop_ds_ =,
	dpin_ds_ =,
	condition_matching_patterns_ds_ =,
	condition_criteria_ds_ =,
	rxdrug_reference_ds_ =,
	PHIN_var_ =,
	subject_id_var_ =,
	output_ds_ =
);

	%put %STR( -> Start of %upcase( &sysmacroname ) macro );
	
	%let DATA_SOURCE_CODE = DPIN;
	
	*split the dataset name into two tokens (the library name and the table name);
	%let input_lib_name = %scan(&dpin_ds_, 1, '.');
	%let input_tbl_name = %scan(&dpin_ds_, 2, '.');
	
	%select_obs_by_study_date(
		input_ds_ = &dpin_ds_,
		output_ds_ = _dpin_1_ds,
		pop_ds_ = &pop_ds_,
		PHIN_var_ = &PHIN_var_,
		subject_id_var_ = &subject_id_var_
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=&input_lib_name nolist;
			delete &input_tbl_name;
		run;
	%end;

	%setup_dpin_for_pattern_matching(
		input_ds_ = _dpin_1_ds,
		output_ds_ = _dpin_2_ds,
		condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
		rxdrug_reference_ds_ = &rxdrug_reference_ds_,
		subject_id_var_ = &subject_id_var_
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=work nolist;
			delete _dpin_1_ds;
		run;
	%end;

	%select_obs_with_conditions(
		input_ds_ = _dpin_2_ds,
		output_ds_ = _dpin_3_ds,
		condition_matching_patterns_ds_ = &condition_matching_patterns_ds_,
		subject_id_var_ = &subject_id_var_,
		data_source_codes_ = ("&DATA_SOURCE_CODE")
	);
	
	%if %eval(&debug_mode_ = 0) %then %do;
		proc datasets lib=work nolist;
			delete _dpin_2_ds;
		run;
	%end;

	%apply_condition_criteria(
		input_ds_ = _dpin_3_ds,
		output_ds_ = &output_ds_,
		data_source_code_ = &DATA_SOURCE_CODE,
		condition_criteria_ds_ = &condition_criteria_ds_,
		subject_id_var_ = &subject_id_var_
	);

	%put %STR( -> End of %upcase( &sysmacroname ) macro );

%mend;

/*
-Determines if matching will be done based on the ATC or DIN coding sytems, or both. Then prepares the DPIN
 dataset to use the correct coding system(s).
*/
%macro setup_dpin_for_pattern_matching(
	input_ds_ =,
	output_ds_ =,
	condition_matching_patterns_ds_ =,
	rxdrug_reference_ds_ =,
	subject_id_var_ =
);

	%put %STR( --> Start of %upcase( &sysmacroname ) macro );
	
	%let dpin_ATC_ds = dpin_ATC_ds;
	%let dpin_DIN_ds = dpin_DIN_ds;
	%let combined_dpin_datasets = combined_dpin_datasets;

	*Determine if matching should be performed using ATC codes, DIN codes, or both;
	proc sql noprint;

		select (count(*) > 0)
		into :match_using_ATC_codes
		from &condition_matching_patterns_ds_
		where data_source = 'DPIN'
		and coding_system = 'ATC';
		
		select (count(*) > 0)
		into :match_using_DIN_codes
		from &condition_matching_patterns_ds_
		where data_source = 'DPIN'
		and coding_system = 'DIN';

	quit;

	*Create DPIN dataset with ATC codes;
	%if (&match_using_ATC_codes = 1) %then %do;

		proc sql noprint;

			create table &dpin_ATC_ds as
			select &subject_id_var_, service_date, 'ATC' as coding_system, upcase(strip(rxref.atc)) as code
			from &input_ds_ as dpin
			inner join &rxdrug_reference_ds_ as rxref
				on dpin.code = rxref.din
			where length(strip(rxref.atc)) > 0;

		quit;

	%end;

	*Create DPIN dataset with DIN codes (DPIN natively uses DIN codes);
	%if (&match_using_DIN_codes = 1) %then %do;
	
		data &dpin_DIN_ds;
		
			set &input_ds_;
			
			keep &subject_id_var_ service_date coding_system code;
			
		run;

	%end;

	%combine_dpin_datasets(
		dpin_ATC_ds_ = &dpin_ATC_ds,
		dpin_DIN_ds_ = &dpin_DIN_ds,
		output_ds_ = &combined_dpin_datasets
	);
	
	*If the same rx drug code occurs more than once for the same person and for the same provided date, the duplicate codes will be removed;
	*This is necessary since multiple DIN codes can map to the same ATC code;
	%if %sysfunc(exist(&combined_dpin_datasets)) %then %do;
		
		proc sort data = &combined_dpin_datasets out = &output_ds_ nodupkey;
			by &subject_id_var_ service_date coding_system code;
		run;
		
	%end;

	proc datasets noprint;
	
		%if %sysfunc(exist(&dpin_ATC_ds)) %then %do;
			delete &dpin_ATC_ds;
		%end;

		%if %sysfunc(exist(&dpin_DIN_ds)) %then %do;
			delete &dpin_DIN_ds;
		%end;
		
	quit;

	%put %STR( --> End of %upcase( &sysmacroname ) macro );
	
%mend;

%macro combine_dpin_datasets(
	dpin_ATC_ds_ =,
	dpin_DIN_ds_ =,
	output_ds_ =
);

	%put %STR( --> Start of %upcase( &sysmacroname ) macro );

	%let ATC_ds_exists = %sysfunc(exist(&dpin_ATC_ds_));
	%let DIN_ds_exists = %sysfunc(exist(&dpin_DIN_ds_));

	data &output_ds_;
	
		%if (&ATC_ds_exists = 1 and &DIN_ds_exists = 1) %then %do;
		
			set &dpin_ATC_ds_ &dpin_DIN_ds_;
			
		%end;
		%else %if (&ATC_ds_exists = 1) %then %do;

			set &dpin_ATC_ds_;

		%end;
		%else %if (&DIN_ds_exists = 1) %then %do;

			set &dpin_DIN_ds_;
			
		%end;
	
	run;
	
	%put %STR( --> End of %upcase( &sysmacroname ) macro );

%mend;

/*
-removes records where the service date does not fall within the study period
*/
%macro select_obs_by_study_date(
	input_ds_ =,
	output_ds_ =,
	pop_ds_ =,
	PHIN_var_ =,
	subject_id_var_ =
);

	%put %STR( --> Start of %upcase( &sysmacroname ) macro );

	/* 
	Select only the records which fall within the study period (start and end dates are included). 
	The PHIN identifier may not be unique within the population dataset (depends upon the study type) so it is 
	important to use the subject_id as the unique identifier.
	*/
	proc sql noprint;

		create table &output_ds_. as
		select pop.&subject_id_var_,
			   pop.study_start_date,
			   pop.study_end_date,
			   input_ds.service_date,
			   input_ds.coding_system,
			   input_ds.code
		from &pop_ds_. as pop
		inner join &input_ds_ as input_ds
			on pop.&PHIN_var_ = input_ds.&PHIN_var_
		where service_date >= study_start_date and
			  service_date <= study_end_date
		order by pop.&subject_id_var_;

	quit;

	%put %STR( --> End of %upcase( &sysmacroname ) macro );
 
%mend;

/*
-select only those diagnostic/prescription codes that indicate a condition that is of interest for a particular project.
-uses a case insensitive search
*/
%macro select_obs_with_conditions(
	input_ds_ =,
	output_ds_ =,
	condition_matching_patterns_ds_ =,
	subject_id_var_ =,
	data_source_codes_ =
);

	%put %STR( --> Start of %upcase( &sysmacroname ) macro );

	/*
	-The like operator in SAS does not work correctly if you just provide a variable/column name (as opposed to providing a string literal).
	 As a result, it is necessary to call strip() or another string function on the column name for pattern matching to work correctly.
	-IMPORTANT: do not remove strip function call in the join clause
	*/
	proc sql noprint;

		create table &output_ds_ as
		select input_ds.&subject_id_var_, 
			   input_ds.service_date, 
			   input_ds.coding_system as ds_coding_system, 
			   input_ds.code, 
			   cond_patterns.condition_id, 
			   cond_patterns.data_source, 
			   cond_patterns.coding_system, 
			   cond_patterns.code_matching_pattern
		from &input_ds_ as input_ds
		inner join &condition_matching_patterns_ds_ as cond_patterns
			on input_ds.coding_system = cond_patterns.coding_system
			and input_ds.code like strip(cond_patterns.code_matching_pattern)
		where cond_patterns.data_source in &data_source_codes_;

	quit;
	
	%put %STR( --> End of %upcase( &sysmacroname ) macro );
	
%mend;

%macro apply_condition_criteria(
	input_ds_ =,
	output_ds_ =,
	data_source_code_ =,
	condition_criteria_ds_ =,
	subject_id_var_ =
);

	%put %STR( --> Start of %upcase( &sysmacroname ) macro );
	
	proc sql noprint;

		*remove duplicate diagnoses for the same condition on the same service date;
		create table &temp_lib_..&data_source_code_._condition_diagnoses_1 as
		select distinct &subject_id_var_, data_source, condition_id, service_date
		from &input_ds_
		order by &subject_id_var_, data_source, condition_id, service_date;

	quit;
	
	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=work nolist;
			delete &input_ds_;
		run;
		
	%end;
	
	proc sql noprint;
	
		*Add time_between_cases value (specific to each condition) by joining the condition criteria lookup table;
		create table &temp_lib_..&data_source_code_._condition_diagnoses_2 as
		select &subject_id_var_, condition_diagnoses.data_source, condition_diagnoses.condition_id, condition_diagnoses.service_date, time_between_cases
		from &temp_lib_..&data_source_code_._condition_diagnoses_1 as condition_diagnoses
		inner join &condition_criteria_ds_ as condition_criteria
			on condition_diagnoses.condition_id = condition_criteria.condition_id
				and condition_diagnoses.data_source = condition_criteria.data_source
		order by &subject_id_var_, data_source, condition_id, service_date;

	quit;

	*Decide which records should be discarded. If the time_between_cases constraint is provided for a disease definition then any cases that are not at least X many
	days apart will be marked to be excluded (keep_record = false);
	data &temp_lib_..&data_source_code_._condition_diagnoses_3;

		set &temp_lib_..&data_source_code_._condition_diagnoses_2;

		*Service dates are to be compared only for the same condition, for the same individual (and for the same data source);
		by &subject_id_var_ data_source condition_id;
		
		*Store the service date from the previous case (for the same condition and same individual) to be compared with the current case;
		retain retaining_service_date_var;

		*All records are kept by default unless time_between_cases is defined and the computed days lapsed is greater than or equal to the time_between_cases constraint.
		This ensures records are kept when disease definitions do not include the time_between_cases constraint (time_between_cases = null);
		keep_record = 1;

		*Save the service_date from the previous case before updating it for next iteration;
		prev_case_service_date = retaining_service_date_var;
		
		if not first.condition_id then do;

			if not missing(time_between_cases) then do;
				
				time_since_last_case = (service_date - prev_case_service_date);
				
				keep_record = (time_since_last_case >= time_between_cases);

			end;

		end;
		else do;
			time_since_last_case = .;
		end;

		*Update the retaining variable with the service_date needed for the next iteration. If the current case is too close (in date/time) to the previous one (keep_record is equal to 0)
		then do not update the retaining variable and leave it to contain the previous service_date value. (This ensures the next service_date will be compared with the service_date from the 
		last valid case.). Also, if it is the last record in the set (same subject and same condition) then set retained date to be null;
		if not last.condition_id and keep_record then do;
			retaining_service_date_var = service_date;
		end;
		else if last.condition_id then do;
			retaining_service_date_var = .;
		end;

		format prev_case_service_date date11.;
		format retaining_service_date_var date11.;
		
	run;

	*Delete the records that were marked to be discarded. These are the records that were too close to the case that came before (time_since_last_case < time_between_cases);
	data &temp_lib_..&data_source_code_._condition_diagnoses_4;
	
		set &temp_lib_..&data_source_code_._condition_diagnoses_3;
	
		if keep_record = 1;
		
		keep &subject_id_var_ data_source condition_id service_date prev_case_service_date time_since_last_case;

	run;

	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._condition_diagnoses_1;
			delete &data_source_code_._condition_diagnoses_2;
			delete &data_source_code_._condition_diagnoses_3;
		run;
		
	%end;

	*Compute the number of cases subjects have of conditions that are NOT time bound. In other words, the minimum number of cases does not need to be met within
	a certain time period (condition_criteria.time_period = NULL);
	proc sql noprint;

		*Add time_period value (specific to each condition) by joining the condition criteria lookup table and only select those 
		records where the condition is defined with time_period = NULL;
		*if the min_num_cases is equal to 1 then these are also considerd non-time-bound diagnoses;
		create table &temp_lib_..&data_source_code_._non_time_bound_diagnoses_1 as 
		select &subject_id_var_, 
			   condition_diagnoses.data_source, 
			   condition_diagnoses.condition_id, 
			   service_date,
			   time_period
		from &temp_lib_..&data_source_code_._condition_diagnoses_4 as condition_diagnoses
		inner join &condition_criteria_ds_ as condition_criteria
			on condition_diagnoses.condition_id = condition_criteria.condition_id
				and condition_diagnoses.data_source = condition_criteria.data_source
		where ( missing(condition_criteria.time_period)
			  or ((not missing(condition_criteria.time_period)) and (min_num_cases eq 1)) );

	quit;

	%if %eval(&debug_mode_ = 1) %then %do;
	
		proc sort data=&temp_lib_..&data_source_code_._non_time_bound_diagnoses_1;
			by &subject_id_var_ data_source condition_id service_date;
		run;
		
	%end;

	proc sql noprint;
	
		*Count the number of cases for each condition for each individual, and also retrieve the date of first occurrence;
		create table &temp_lib_..&data_source_code_._non_time_bound_diagnoses_2 as
		select &subject_id_var_, 
			   data_source, 
			   condition_id, 
			   min(service_date) format date11. as first_occurrence, 
			   count(*) as num_cases
		from &temp_lib_..&data_source_code_._non_time_bound_diagnoses_1 as condition_diagnoses
		group by &subject_id_var_, 
			     data_source, 
			     condition_id;

		*Filter out any conditions where the minimum number of cases constraint was not met;
		create table &temp_lib_..&data_source_code_._non_time_bound_diag_final as
		select &subject_id_var_, 
			   condition_diagnoses.data_source, 
			   condition_diagnoses.condition_id,
			   first_occurrence,
			   num_cases
		from &temp_lib_..&data_source_code_._non_time_bound_diagnoses_2 as condition_diagnoses
		inner join &condition_criteria_ds_ as condition_criteria
			on condition_diagnoses.condition_id = condition_criteria.condition_id
				and condition_diagnoses.data_source = condition_criteria.data_source
		where num_cases >= min_num_cases;

	quit;
	
	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._non_time_bound_diagnoses_1;
			delete &data_source_code_._non_time_bound_diagnoses_2;
		run;
		
	%end;

	*Compute the number of cases subjects have of conditions that ARE time bound. In other words, the minimum number of cases MUST occur within
	a certain time period (condition_criteria.time_period != NULL);
	proc sql noprint;
	
		*Add time_period value (specific to each condition) by joining the condition criteria lookup table and only select those 
		records where the condition is defined with time_period != NULL. Also, calculate the time period end date;
		*Note: it does not make sense to look for cases to be in a certain time period if there is only one required case, therefore min_num_cases must be greater than 1;
		create table &temp_lib_..&data_source_code_._time_bound_diagnoses_1 as
		select &subject_id_var_, 
			   condition_diagnoses.data_source, 
			   condition_diagnoses.condition_id, 
			   service_date,
			   time_period,
			   (service_date + time_period) as time_period_end_date format date11.
		from &temp_lib_..&data_source_code_._condition_diagnoses_4 as condition_diagnoses
		inner join &condition_criteria_ds_ as condition_criteria
			on condition_diagnoses.condition_id = condition_criteria.condition_id
				and condition_diagnoses.data_source = condition_criteria.data_source
		where ((not missing(condition_criteria.time_period)) and (min_num_cases > 1));

	quit;
	
	%if %eval(&debug_mode_ = 1) %then %do;
	
		proc sort data=&temp_lib_..&data_source_code_._time_bound_diagnoses_1;
			by &subject_id_var_ data_source condition_id service_date;
		run;
		
	%end;

	proc sql noprint;
	
		*Join the table containing the time bound diagnoses with itself (&data_source_code_._time_bound_diagnoses_1) in order to count how many cases occurred with each time period.
		A separate time window is created for each case of a particular condition (for a particular individual). The start date of each of these time windows is the 
		service date where the condition was diagnosed;
		create table &temp_lib_..&data_source_code_._time_bound_diagnoses_2 as
		select condition_diagnoses_1.&subject_id_var_, 
			   condition_diagnoses_1.data_source, 
			   condition_diagnoses_1.condition_id, 
			   condition_diagnoses_1.service_date as time_period_start_date,
			   condition_diagnoses_1.time_period_end_date,
			   condition_diagnoses_2.service_date as service_date
		from &temp_lib_..&data_source_code_._time_bound_diagnoses_1 as condition_diagnoses_1
		left outer join &temp_lib_..&data_source_code_._time_bound_diagnoses_1 as condition_diagnoses_2
			on condition_diagnoses_1.&subject_id_var_ = condition_diagnoses_2.&subject_id_var_
			and condition_diagnoses_1.data_source = condition_diagnoses_2.data_source 
			and condition_diagnoses_1.condition_id = condition_diagnoses_2.condition_id
		where (condition_diagnoses_2.service_date >= condition_diagnoses_1.service_date 
			   and condition_diagnoses_2.service_date < condition_diagnoses_1.time_period_end_date);

	quit;

	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._time_bound_diagnoses_1;
		run;
		
	%end;

	%if %eval(&debug_mode_ = 1) %then %do;
	
		proc sort data=&temp_lib_..&data_source_code_._time_bound_diagnoses_2;
			by &subject_id_var_ data_source condition_id time_period_start_date service_date;
		run;
		
	%end;
	
	proc sql noprint;
	
		*For each time period, calculate the date of first occurence and the number of cases;
		create table &temp_lib_..&data_source_code_._time_bound_diagnoses_3 as
		select &subject_id_var_, 
			   data_source, 
			   condition_id, 
			   time_period_start_date,
			   min(service_date) format date11. as first_occurrence_in_time_period, 
			   count(*) as num_cases_in_time_period
		from &temp_lib_..&data_source_code_._time_bound_diagnoses_2
		group by &subject_id_var_, 
			     data_source, 
			     condition_id, 
			     time_period_start_date;

	quit;

	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._time_bound_diagnoses_2;
		run;
		
	%end;

	proc sql noprint;
	
		*Filter results to only contain time periods where the minumum number of cases was met;
		*Note: it does not make sense to look for cases to be in a certain time period if there is only one required case, therefore min_num_cases must be greater than 1 (conditions
			   not meeting this are excluded);
		create table &temp_lib_..&data_source_code_._time_bound_diagnoses_4 as
		select &subject_id_var_, 
			   condition_diagnoses.data_source, 
			   condition_diagnoses.condition_id,
			   time_period_start_date,
			   first_occurrence_in_time_period,
			   num_cases_in_time_period
		from &temp_lib_..&data_source_code_._time_bound_diagnoses_3 as condition_diagnoses
		inner join &condition_criteria_ds_ as condition_criteria
			on condition_diagnoses.condition_id = condition_criteria.condition_id
				and condition_diagnoses.data_source = condition_criteria.data_source
		where ((not missing(condition_criteria.time_period)) and (min_num_cases > 1))
		and num_cases_in_time_period >= min_num_cases;

	quit;

	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._time_bound_diagnoses_3;
		run;
		
	%end;

	proc sql noprint;

		*Summarize data to provide the maximum number of cases (within a time period) and the date of first occurence for each condition, per data source, that an individual has;
		create table &temp_lib_..&data_source_code_._time_bound_diagnoses_final as
		select &subject_id_var_, 
			   data_source, 
			   condition_id,
			   min(first_occurrence_in_time_period) format date11. as first_occurrence, 
			   max(num_cases_in_time_period) as max_num_cases_in_time_period
		from &temp_lib_..&data_source_code_._time_bound_diagnoses_4
		group by &subject_id_var_, 
			     data_source, 
			     condition_id;

	quit;
	
	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._time_bound_diagnoses_4;
		run;
		
	%end;
	
	proc sql noprint;
	
		*Combine the non-time bound diagnoses as well as those that must occur within a specified time period;
		*use union rather than union all to remove duplicates;
		create table &temp_lib_..&data_source_code_._all_condition_diagnoses as
		select &subject_id_var_, data_source, condition_id, first_occurrence
		from &temp_lib_..&data_source_code_._non_time_bound_diag_final
		union
		select &subject_id_var_, data_source, condition_id, first_occurrence
		from &temp_lib_..&data_source_code_._time_bound_diagnoses_final;
	
	quit;

	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._non_time_bound_diag_final;
			delete &data_source_code_._time_bound_diagnoses_final;
		run;
		
	%end;

	proc sql noprint;

		*The resultset may have two occurences for the same disease for the same person (for example, one from the physician diagnostic code (PHYS) and one from the physician tariff code (TARF)). 
		Therefore, use aggregation to select only unique condition outcomes and the earliest date of first occurence;
		create table &output_ds_ as
		select &subject_id_var_, condition_id, min(first_occurrence) format date11. as first_occurrence
		from &temp_lib_..&data_source_code_._all_condition_diagnoses
		group by &subject_id_var_, condition_id;

	quit;

	%if %eval(&debug_mode_ = 0) %then %do;
	
		proc datasets lib=&temp_lib_ nolist;
			delete &data_source_code_._all_condition_diagnoses;
			delete &data_source_code_._condition_diagnoses_4;
		run;
		
	%end;

	%put %STR( --> End of %upcase( &sysmacroname ) macro );

%mend;

%macro create_final_ds(
	pop_ds_ =,
	hosp_condition_cases_ =,
	phys_condition_cases_ =,
	dpin_condition_cases_ =,
	conditions_ds_ =,
	subject_id_var_ =,
	output_widefmt_binary_ds_ =,
	output_widefmt_date_ds_ =,
	output_longfmt_ds_ =
);

	%put %STR( --> Start of %upcase( &sysmacroname ) macro );

	%local pop_library_name;
	%local pop_table_name;

	*split the population dataset name into two tokens (the library name and the table name);
	%let pop_library_name = %scan(&pop_ds_, 1, '.');
	%let pop_table_name = %scan(&pop_ds_, 2, '.');

	%if &output_widefmt_binary_ds_ = . %then %do;
		%let output_widefmt_binary_ds_ = work.conditions_widefmt;
	%end;
	
	%if &output_widefmt_date_ds_ = . %then %do;
		%let output_widefmt_date_ds_ = work.conditions_widefmt_date;
	%end;
	
	%if &output_longfmt_ds_ = . %then %do;
		%let output_longfmt_ds_ = work.conditions_longfmt;
	%end;

	*Combine datasets;
	proc sql noprint;

		select length
		into :subject_id_length
		from sashelp.vcolumn
		where libname = %UPCASE("&pop_library_name")
		and memname = %UPCASE("&pop_table_name")
		and name = "&subject_id_var_";
		
		create table work.all_cases
		(
			&subject_id_var_ varchar(&subject_id_length),
			condition_id varchar(32),
			first_occurrence int format date11.
		);
	
	quit;
	
	%if %sysfunc(exist(&hosp_condition_cases_)) %then %do;
	
		proc sql noprint;
		
			insert into work.all_cases
			select *
			from &hosp_condition_cases_;
	
		quit;
		
	%end;
	
	%if %sysfunc(exist(&phys_condition_cases_)) %then %do;

		proc sql noprint;
		
			insert into work.all_cases
			select *
			from &phys_condition_cases_;
	
		quit;
		
	%end;

	%if %sysfunc(exist(&dpin_condition_cases_)) %then %do;

		proc sql noprint;

			insert into work.all_cases
			select *
			from &dpin_condition_cases_;

		quit;

	%end;

	*Remove duplicate cases (i.e. the same individual might be recorded with the same condition in more than one data source such as hosp, phys, and dpin);
	proc sql noprint;

		create table &output_longfmt_ds_ as
		select &subject_id_var_ label='Subject ID', condition_id, min(first_occurrence) format date11. as first_occurrence
		from work.all_cases
		group by &subject_id_var_, condition_id;

	quit;

	%let input_ds = &output_longfmt_ds_;
	%let output_binary_ds = work.conditions_1_binary;
	%let output_date_ds = work.conditions_1_date;

	*Prepare the datasets to be transposed to wide format: add the condition name, add a boolean flag, and sort the data accordingly;
	proc sql noprint;

		create table &output_binary_ds as 
		select &subject_id_var_, outcomes.condition_id, condition_name, 1 as has_condition
		from &input_ds as outcomes
		inner join &conditions_ds_ as conditions
			on outcomes.condition_id = conditions.condition_id
		order by &subject_id_var_, outcomes.condition_id;
		
		create table &output_date_ds as 
		select &subject_id_var_, outcomes.condition_id, condition_name, first_occurrence
		from &input_ds as outcomes
		inner join &conditions_ds_ as conditions
			on outcomes.condition_id = conditions.condition_id
		order by &subject_id_var_, outcomes.condition_id;

	quit;

	%let input_binary_ds = work.conditions_1_binary;
	%let input_date_ds = work.conditions_1_date;
	%let output_binary_ds = work.conditions_2_binary;
	%let output_date_ds = work.conditions_2_date;

	proc transpose data=&input_binary_ds out=&output_binary_ds (drop = _NAME_);
		by &subject_id_var_;
		id condition_id;
		idlabel condition_name;
		var has_condition;
	run;

	proc transpose data=&input_date_ds out=&output_date_ds (drop = _NAME_);
		by &subject_id_var_;
		id condition_id;
		idlabel condition_name;
		var first_occurrence;
	run;

	%let input_binary_ds = work.conditions_2_binary;
	%let input_date_ds = work.conditions_2_date;
	%let output_binary_ds = work.conditions_3_binary;
	%let output_date_ds = &output_widefmt_date_ds_;

	*Add subjects that do not have any of the listed conditions;
	proc sql noprint;

		create table &output_binary_ds as
		select *
		from &input_binary_ds;

		insert into &output_binary_ds (&subject_id_var_)
		select &subject_id_var_
		from &pop_ds_
		where &subject_id_var_ not in (select &subject_id_var_ from &input_binary_ds);

	quit;
	
	proc sql noprint;

		create table &output_date_ds as
		select *
		from &input_date_ds;

		insert into &output_date_ds (&subject_id_var_)
		select &subject_id_var_
		from &pop_ds_
		where &subject_id_var_ not in (select &subject_id_var_ from &input_date_ds);

	quit;

	%let input_binary_ds = work.conditions_3_binary;
	%let output_binary_ds = &output_widefmt_binary_ds_;

	*Wherever values are missing, set those fields to zero;
	proc stdize data=&input_binary_ds out=&output_binary_ds reponly missing=0;
	run;

	%put %STR( --> End of %upcase( &sysmacroname ) macro );

%mend;
