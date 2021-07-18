import {getText} from './translation.js';
const slugify = str => str.toLowerCase().trim().replace(/[^\w\s-]/g, '').replace(/[\s_-]+/g, '-').replace(/^-+|-+$/g, '');

// TODO: Show max time and memory too in summary

export class SubmissionManager {
	constructor(id, replace_id, lang) {
		this.id = id
		this.replace_id = replace_id
		this.lang = lang || "en"
		this.poll_mu = false
		this.subAuthor = {}
		this.subProblem = {}
		this.problemEditor = false
		this.sub = {}
		this.subTests = []
		this.poller = null
		this.finished = false

		this.subTasks = []
		this.subTestIDs = {}
	}

	async startPoller() {
		console.log("Started poller")
		await this.poll()
		if(!this.finished) {
			this.poller = setInterval(async () => this.poll(), 2000)
		}
	}

	stopPoller() {
		if(this.poller == null) {
			return
		}
		console.log("Stopped poller")
		clearInterval(this.poller)
		this.poller = null
	}

	downloadCode() {
		var file = new Blob([this.sub.code], {type: 'text/plain'})
		var filename = `${slugify(this.subProblem.name, {lower: true})}-${this.id}.${this.sub.language}`
		bundled.downloadBlob(file, filename);
	}

	async copyCode() {
		await navigator.clipboard.writeText(this.sub.code)
		bundled.createToast({status: "success", description: this.getText("copied")})
	}

	async poll() {
		if(this.poll_mu === false) this.poll_mu = true
		else return
		console.log("Poll", this.id)
		let res = await bundled.getCall("/submissions/getByID", {id: this.id, expanded: true})
		if(res.status !== "success") {
			bundled.apiToast(res)
			console.error(res)
			this.poll_mu = false
			return
		}

		console.log(res)

		res = res.data
		if(res.subtests) {
			this.subTests = res.subtests
			this.subTestIDs = {}
			for(let subtest of res.subtests) {
				this.subTestIDs[subtest.pb_test.id] = subtest;
			}
		}
		
		this.sub = res.sub
		this.subEditor = res.sub_editor
		this.problemEditor = res.problem_editor
		this.subAuthor = res.author
		this.subProblem = res.problem
		if(res.subtasks) {
			this.subTasks = res.subtasks
		}

		if(this.sub.status === "finished") {
			this.stopPoller()
			this.finished = true
		}

		this.render()
		this.poll_mu = false
	}

	async toggleVisible() {
		let res = await bundled.postCall("/submissions/setVisible", {visible: !this.sub.visible, id: this.id});
		bundled.apiToast(res)
		this.sub.visible = !this.sub.visible
		this.render();
	}

	async toggleQuality() {
		let res = await bundled.postCall("/submissions/setQuality", {quality: !this.sub.quality, id: this.id});
		this.sub.quality = !this.sub.quality
		bundled.apiToast(res)
		this.render();
	}

	summaryNode() {
		let rez = document.createElement('div')
		let html = ""
		html += `<h2>${this.getText("sub")} ${this.sub.id}</h2>` +
			`<p>${this.getText("author")}: <a href="/profile/${this.subAuthor.name}">${this.subAuthor.name}</a></p>` +
			`<p>${this.getText("problem")}: <a href="/problems/${this.subProblem.id}">${this.subProblem.name}</a></p>` +
			`<p>${this.getText("uploadDate")}: ${bundled.parseTime(this.sub.created_at)}</p>` +
			`<p>${this.getText("status")}: ${this.sub.status}</p>` +
			`<p>${this.getText("lang")}: ${this.sub.language}</p>` 
		if(this.sub.quality) {
			html += `<p><i class="fas fa-star text-yellow-300"></i> ${this.getText("qualitySub")}</p>`
		}
		if(this.sub.code) {
			html += `<p>${this.getText("size")}: ${bundled.sizeFormatter(this.sub.code.length)}</p>`
		}
		if(this.subProblem.default_points > 0) {
			html += `<p>${this.getText("defaultPoints")}: ${this.subProblem.default_points}</p>`
		}
		if(this.sub.status === 'finished') {
			html += `<p>${this.getText("score")}: ${this.sub.score}</p>`
		}
		if(this.sub.compile_error.Bool) {
			html += `<h4>${this.getText("compileErr")}</h4><h5>${this.getText("compileMsg")}:</h5><pre>${this.sub.compile_message.String}</pre>`
		}
		rez.innerHTML = html
		return rez;
	}

	tableColGen(text) {
		let td = document.createElement('td')
		td.innerHTML = text
		return td
	}

	subTasksNode() {
		let rezz = document.createElement('div');
		rezz.classList.add('my-2');
		let rez = document.createElement('div')
		rez.classList.add('list-group', 'my-1', 'list-group-mini')
		for(let subtask of this.subTasks) {
			let row = document.createElement('details')
			row.classList.add('list-group-item')
			
			let sum = document.createElement('summary')
			sum.classList.add('flex', 'justify-between')
			
			let stk_score = 100;
			let subtests = document.createElement('div')
			subtests.classList.add('list-group', 'm-1', 'list-group-mini')
			for(let testID of subtask.tests) {
				let roww = document.createElement('div')
				roww.classList.add('list-group-item', 'flex', 'justify-between')

				let actualTest = this.subTestIDs[testID];
				if(actualTest.subtest.score < stk_score) {
					stk_score = actualTest.subtest.score;
				}
				roww.innerHTML = `<span>${this.getText("test")} #${actualTest.pb_test.visible_id}</span><span class="rounded-full py-1 px-2 text-base text-white font-semibold" style="background-color: ${bundled.getGradient(actualTest.subtest.score, 100)}">${Math.round(subtask.score * actualTest.subtest.score / 100.0)} / ${subtask.score}</span>`
				
				subtests.appendChild(roww)
			}

			sum.innerHTML = `<span>${this.getText("subTask")} #${subtask.visible_id}</span><span class="rounded-full py-1 px-2 text-base text-white font-semibold" style="background-color: ${bundled.getGradient(stk_score, 100)}">${Math.round(subtask.score * stk_score / 100.0)} / ${subtask.score}</span>`
			
			row.appendChild(sum)
			row.appendChild(subtests)
			rez.appendChild(row)
		}
		rezz.appendChild(rez)

		let tmp = document.createElement('details');
		let tmp1 = document.createElement('summary');
		tmp1.innerText = this.getText("seeTests");
		tmp.appendChild(tmp1);
		tmp.appendChild(this.tableNode());
		rezz.appendChild(tmp);

		return rezz
	}

	tableNode() {
		let rez = document.createElement('table')
		rez.classList.add('kn-table')
		let head = document.createElement('thead')
		head.innerHTML = `<tr><th class="py-1" scope="col">${this.getText("id")}</th><th scope="col">${this.getText("time")}</th><th scope="col">${this.getText("memory")}</th><th scope="col">${this.getText("verdict")}</th><th scope="col">${this.getText("score")}</th>${this.problemEditor ? `<th scope='col'>${this.getText("output")}</th>` : ""}${this.subTasks.length > 0 ? `<th scope='col'>${this.getText("subTasks")}</th>` : ""}</tr>`
		let body = document.createElement('tbody')
		for(let test of this.subTests) {
			let row = document.createElement('tr')
			row.classList.add('kn-table-row')
			
			let vid = document.createElement('th')
			vid.innerText = test.pb_test.visible_id
			vid.classList.add('py-2')
			vid.scope = "row"
			row.appendChild(vid)
			
			let time = this.tableColGen("")
			let mem = this.tableColGen("")
			let verdict = this.tableColGen(`<div class='fas fa-spinner animate-spin' role='status'></div> ${this.getText("waiting")}`)
			let score = this.tableColGen(`${Math.round(test.pb_test.score * test.subtest.score / 100.0)} / ${test.pb_test.score}`)
			if(this.subTasks.length > 0) {
				score.innerHTML = `${test.subtest.score}% ${this.getText("correct")}`
			}
			if(test.subtest.done) {
				verdict.innerHTML = test.subtest.verdict
				
				time.innerHTML = Math.floor(test.subtest.time * 1000) + " ms";
				mem.innerHTML = bundled.sizeFormatter(test.subtest.memory*1024, 1, true)

				score.classList.add("text-black")
				score.style = "background-color:" + bundled.getGradient(test.subtest.score, 100) + ";"
			} else {
				score.innerHTML = "-"
			}

			row.appendChild(time)
			row.appendChild(mem)
			row.appendChild(verdict)
			row.appendChild(score)
			if(this.problemEditor) {
				let out = this.tableColGen("")
				if(test.subtest.done) {
					out.innerHTML = `<a href="/proposer/get/subtest_output/${test.subtest.id}">${this.getText("output")}</a>`
				}
				row.appendChild(out)
			}
			if(this.subTasks.length > 0) {
				let out = this.tableColGen(this.testSubTasks(test.pb_test.id).join(', '))
				row.appendChild(out)
			}

			body.appendChild(row);
		}
		rez.appendChild(head)
		rez.appendChild(body)
		return rez;
	}

	testSubTasks(test_id) {
		let stks = [];
		for(let st of this.subTasks) {
			if(st.tests.includes(test_id)) {
				stks.push(st.visible_id);
			}
		}
		return stks;
	}

	codeNode() {
		let rez = document.createElement('div')
		
		// header
		let header = document.createElement('h3')
		header.innerText = this.getText("source")
		rez.appendChild(header)

		// code
		let code = document.createElement('pre')
		let c = document.createElement('code')
		c.classList.add('hljs', this.sub.language)
		c.innerHTML = hljs.highlight(this.sub.language, this.sub.code).value
		code.appendChild(c)
		rez.appendChild(code)

		let dv = document.createElement('div')
		dv.classList.add('block', 'my-2')

		let btn = document.createElement('button')
		btn.classList.add('btn', 'btn-blue', 'mr-2', 'text-semibold', 'text-lg')
		btn.innerText = this.getText("copy")
		btn.onclick = async () => await this.copyCode()
		dv.appendChild(btn)

		let btn1 = document.createElement('button')
		btn1.classList.add('btn', 'btn-blue', 'text-semibold', 'text-lg')
		btn1.innerText = this.getText("download")
		btn1.onclick = () => this.downloadCode()
		dv.appendChild(btn1)
		rez.appendChild(dv)

		if(this.subEditor) {
			let btn = document.createElement('button');
			btn.classList.add('btn', 'btn-blue', 'block', 'my-2', 'text-semibold', 'text-lg');
			btn.innerHTML = `<i class="fas fa-share-square mr-2"></i>${this.getText("makeCode")} ${this.sub.visible ? this.getText("invisible") : this.getText("visible")}</button>`;
			btn.onclick = () => this.toggleVisible();
			rez.appendChild(btn);
		}
		if(this.problemEditor) {
			let btn = document.createElement('button');
			btn.classList.add('btn', 'btn-blue', 'block', 'my-2', 'text-semibold', 'text-lg');
			btn.innerHTML = `<i class="fas fa-star mr-2"></i>${this.sub.quality ? this.getText("makeQuality") : this.getText("dropQuality")}</button>`;
			btn.onclick = () => this.toggleQuality();
			rez.appendChild(btn);
		}

		return rez;
	}

	viewNode() {
		let rez = document.createElement('div')
		rez.appendChild(this.summaryNode())
		if(this.subTests.length > 0 && !this.sub.compile_error.Bool) {
			if(this.subTasks.length > 0) {
				rez.appendChild(this.subTasksNode())
			} else {
				rez.appendChild(this.tableNode())
			}
		}
		if(this.sub.code != null) {
			rez.appendChild(this.codeNode())
		}
		return rez;
	}

	render() {
		let node = this.viewNode()
		node.id = this.replace_id

		let target = document.getElementById(this.replace_id);
		target.parentNode.replaceChild(node, target);
	}

	getText(key) {
		return getText(this.lang, key)
	}
}