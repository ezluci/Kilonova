import { fromBase64 } from "js-base64";
import { h, Fragment } from "preact";
import register from "preact-custom-element";
import { useEffect, useState } from "preact/hooks";
import { getCall } from "../api/client";
import { apiToast } from "../toast";
import getText from "../translation";
import { BigSpinner, formatScoreStr } from "./common";

type Sublist = {
	id: number;
	title: string;
	author_id: number;
	num_problems: number;
};

type FullList = {
	id: number;
	title: string;
	description: string;
	author_id: number;
	num_problems: number;
	list: number[];
	sublists: Sublist[];
};

type ScoredProblem = Problem & {
	is_editor: boolean;
	score_user_id?: number | null;
	max_score?: number | null;
};

function isProblemEditor(pb: ScoredProblem): boolean {
	if (window.platform_info.admin) {
		return true;
	}
	return pb.is_editor;
}

type ProblemScore = { [problem: number]: number };
type SublistSolved = { [num_solved: number]: number };

function SPBMaxScore({ pb }: { pb: ScoredProblem }) {
	console.log(pb);
	if (typeof pb.score_user_id == "undefined" || pb.score_user_id == null) {
		return <></>;
	}
	if (typeof pb.max_score == "undefined" || pb.max_score == null || pb.max_score < 0) {
		return <>-</>;
	}
	if (pb.scoring_strategy == "acm-icpc") {
		if (Math.abs(pb.max_score - 100) < 0.01) {
			return <i class="fas fa-fw fa-check"></i>;
		}
		return <i class="fas fa-fw fa-xmark"></i>;
	}
	return <>{formatScoreStr(pb.max_score.toFixed(pb.score_precision))}</>;
}

export function Problems({ pbs, listID }: { pbs: ScoredProblem[]; listID?: number }) {
	return (
		<div class="list-group">
			{pbs.map((pb) => (
				<a
					href={`/problems/${pb.id}${typeof listID !== "undefined" ? `?list_id=${listID}` : ""}`}
					class="list-group-item flex justify-between"
					key={pb.id}
				>
					<span>
						{pb.name} (#{pb.id})
					</span>
					{window.platform_info.user_id > 0 && (
						<div>
							{isProblemEditor(pb) &&
								(pb.visible ? (
									<span class="badge badge-green">{getText("published")}</span>
								) : (
									<span class="badge badge-red">{getText("unpublished")}</span>
								))}{" "}
							<span class="badge">
								<SPBMaxScore pb={pb} />
							</span>
						</div>
					)}
				</a>
			))}
		</div>
	);
}

type QueryResult = {
	list: FullList;
	numSolved: number;
	description: string;
	problems: ScoredProblem[];
	problemScores: ProblemScore;
	numSubSolved: SublistSolved;
};

export function Sublist({ list, numsolved }: { list: Sublist; numsolved?: number }) {
	let [loading, setLoading] = useState(false);
	let [expanded, setExpanded] = useState(false);
	let [fullData, setFullData] = useState<FullList | undefined>(undefined);
	let [numSolved, setNumSolved] = useState<number>(numsolved ?? -1);
	let [descHTML, setDescHTML] = useState("");
	let [problems, setProblems] = useState<ScoredProblem[]>([]);
	let [sublistSolved, setSublistSolved] = useState<SublistSolved>({});

	async function load() {
		setLoading(true);
		const res = await getCall<QueryResult>(`/problemList/${list.id}/complex`, {});
		if (res.status === "error") {
			apiToast(res);
			return;
		}
		setFullData(res.data.list);
		setNumSolved(res.data.numSolved);
		setDescHTML(res.data.description);
		setProblems(res.data.problems);
		setSublistSolved(res.data.numSubSolved);
		setLoading(false);
	}

	useEffect(() => {
		if (expanded && fullData == undefined && !loading) {
			load();
		}
	}, [expanded]);

	return (
		<details class="list-group-head" onToggle={(e) => setExpanded(e.currentTarget.open)}>
			<summary class="pb-1 mt-1">
				<span>
					{list.title} <a href={`/problem_lists/${list.id}`}>(#{list.id})</a>
				</span>
				{list.num_problems > 0 &&
					((numSolved >= 0 && <span class="float-right badge">{getText("num_solved", numSolved, list.num_problems)}</span>) || (
						<span class="float-right badge">{list.num_problems == 1 ? getText("single_problem") : getText("num_problems", list.num_problems)}</span>
					))}
			</summary>
			{loading && <BigSpinner />}
			{fullData && (
				<>
					{descHTML && (
						<div class="list-group mt-2">
							<div class="list-group-head" dangerouslySetInnerHTML={{ __html: descHTML }}></div>
						</div>
					)}
					{fullData.sublists.length > 0 && (
						<div class="list-group mt-2">
							{fullData.sublists.map((val) => (
								<Sublist
									list={val}
									numsolved={Object.keys(sublistSolved).includes(val.id.toString()) ? sublistSolved[val.id] : undefined}
									key={val.id.toString() + "_" + list.id.toString()}
								/>
							))}
						</div>
					)}
					{problems.length > 0 && (
						<div class="mt-2">
							<Problems pbs={problems} listID={fullData.id} />
						</div>
					)}
				</>
			)}
		</details>
	);
}

export function DOMSublist({ encoded, numsolved }: { encoded: string; numsolved: string }) {
	let numSolved: number | undefined = parseInt(numsolved);
	if (isNaN(numSolved)) {
		numSolved = undefined;
	}

	return <Sublist list={JSON.parse(fromBase64(encoded))} numsolved={numSolved}></Sublist>;
}

register(DOMSublist, "kn-dom-sublist", ["encoded", "numsolved"]);
