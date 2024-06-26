export {};

declare global {
	// Base types

	type UserBrief = {
		id: number;
		name: string;
		admin: boolean;
		proposer: boolean;

		display_name: string;

		generated: boolean;
	};

	type Problem = {
		id: number;
		created_at: string;
		name: string;
		test_name: string;
		default_points: number;
		visible: boolean;
		time_limit: number;
		memory_limit: number;
		source_credits: string;
		source_size: number;
		console_input: boolean;
		scoring_strategy: "sum_subtasks" | "max_submission" | "acm-icpc";
		score_precision: number;
		published_at?: string;

		score_scale: number;
	};

	type ShallowProblemList = {
		id: number;
		title: string;
		author_id: number;
		sidebar_hidable: boolean;
		featured_checklist: boolean;
		num_problems: number;
	};

	type ProblemList = {
		id: number;
		created_at: string;
		author_id: number;
		title: string;
		description: string;
		list: number[];
		num_problems: number;
		sidebar_hidable: boolean;
		featured_checklist: boolean;
		sublists: ShallowProblemList[];
	};

	type Submission = {
		id: number;
		created_at: string;
		user_id: number;
		problem_id: number;
		language: string;
		code_size: number;
		compile_error: boolean;
		compile_message?: string;
		contest_id: number | null;
		max_time: number;
		max_memory: number;
		score: number;
		status: string;
		score_precision: number;

		compile_time: number | null;

		submission_type: "classic" | "acm-icpc";
		icpc_verdict: string | null;
	};
	type SubTest = {
		id: number;
		created_at: string;
		done: boolean;
		skipped: boolean;
		verdict: string;
		time: number;
		memory: number;
		percentage: number;
		test_id?: number;
		submission_id: number;

		visible_id: number;
		score: number;
	};

	type SubmissionSubTask = {
		id: number;
		created_at: string;

		submission_id: number;
		user_id: number;
		subtask_id?: number;

		problem_id: number;
		visible_id: number;
		score: number;
		final_percentage?: number;

		subtests: number[];
	};

	// Derived types

	type FullSubmission = Submission & {
		author: UserBrief;
		problem: Problem;
		subtests: SubTest[];
		subtasks: SubmissionSubTask[];

		code: string;

		problem_editor: boolean;
		truly_visible: boolean;
	};

	// Contest types
	type Question = {
		id: number;
		asked_at: string;
		responded_at?: string;
		text: string;
		response?: string;
		author_id: number;
		contest_id: number;
	};

	type Announcement = {
		id: number;
		created_at: string;
		contest_id: number;
		text: string;
	};
}
