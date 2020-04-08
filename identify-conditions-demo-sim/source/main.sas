* Purpose: showcase of identifying conditions using simulated data;

/************************************************************ 
	Setup the working root first;
*************************************************************/
%let WORKING_ROOT = /folders/myfolders/conditions-demo;
*************************************************************;

%let SOURCE_PATH  = &WORKING_ROOT/source;
%let DATA_PATH    = &WORKING_ROOT/data;
%let CONFIG_PATH  = &WORKING_ROOT/config;
%let RESULT_PATH  = &WORKING_ROOT/results;

%let hosp_file = "&DATA_PATH/hosp.csv";
%let phys_file = "&DATA_PATH/clinic.csv";
%let pop_file  = "&DATA_PATH/person.csv";
%let dpin_file = "&DATA_PATH/dpin.csv";

%let cond_file  = "&CONFIG_PATH/conditions.csv";
%let crit_file  = "&CONFIG_PATH/criteria.csv";
%let patt_file  = "&CONFIG_PATH/patterns.csv";
%let rxref_file = "&CONFIG_PATH/rx_reference.csv";

%let output_binary = "&RESULT_PATH/conditions_binary.csv";
%let output_date   = "&RESULT_PATH/conditions_date.csv";
%let output_long   = "&RESULT_PATH/conditions_long.csv";
%let output_full   = "&RESULT_PATH/conditions_full.csv";

* Run the conditions macro;
* study period start date: Jan 01, 1995;
* study period data:       Dec 31, 2004;
* inclusion criteria:      40 years old;
* target condion:          diabetes CVD;

%let study_start = '01jan1995'd;
%let study_end   = '31dec2004'd;

%let age_lb = 40;

%include "/folders/myfolders/conditions-demo/source/preprocess_data.sas";

%include "/folders/myfolders/conditions-demo/source/identify_conditions.sas";

%identify_conditions(
	pop_ds_                         = work.study_pop_valid,
	hosp_ds_                        = work.processed_hosp,
	phys_ds_                        = work.processed_phys,
	dpin_ds_                        = work.processed_dpin,
	conditions_ds_                  = work.diab_conditions,
	condition_matching_patterns_ds_ = work.diab_patterns,
	condition_criteria_ds_          = work.diab_criteria,
	rxdrug_reference_ds_            = work.rxref,
	PHIN_var_                       = subject_id,
	subject_id_var_                 = subject_id,
	output_widefmt_binary_ds_       = conds_cases_binary,
	output_widefmt_date_ds_         = conds_cases_date,
	output_longfmt_ds_              = conds_cases_long,
	debug_mode_                     = 0
);
	
* Clean up intermediate data;
proc datasets library = work nolist;
	delete all_cases
		   combined_dpin:
		   conditions_:
		   diab_conditions
		   diab_criteria
		   diab_patterns
		   dpin_:
		   hosp_:
		   phys_:
		   rxref;
quit;

proc sql;
	create table final_report as
	select pop.*,
		   conds.chron_diab as chron_diab_date,
		   conds.chron_card as chron_card_date
	from work.study_pop as pop
	left outer join work.conds_cases_date as conds
	on pop.subject_id = conds.subject_id;
quit;

data final_report;
	set final_report;
	drop study_start_date
		 study_end_date;
run;

* Output results to csv files;
/* proc export  */
/* 	data      = work.conds_cases_binary */
/*     outfile   = &output_binary */
/*     dbms      = dlm */
/*     replace; */
/*     delimiter = ','; */
/* run; */
/*  */
/* proc export  */
/* 	data      = work.conds_cases_date */
/*     outfile   = &output_date */
/*     dbms      = dlm */
/*     replace; */
/*     delimiter = ','; */
/* run; */
/*  */
/* proc export  */
/* 	data      = work.conds_cases_long */
/*     outfile   = &output_long */
/*     dbms      = dlm */
/*     replace; */
/*     delimiter = ','; */
/* run; */

proc export 
	data      = work.final_report
    outfile   = &output_full
    dbms      = dlm
    replace;
    delimiter = ',';
run;