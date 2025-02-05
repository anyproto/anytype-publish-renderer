const katex = require('katex');

class UtilCommon {

	fromHtmlSpecialChars (s: string) {
		return String(s || '').replace(/(&lt;|&gt;|&amp;)/g, (s: string, p: string) => {
			if (p == '&lt;') p = '<';
			if (p == '&gt;') p = '>';
			if (p == '&amp;') p = '&';
			return p;
		});
	};

	getLatex (html: string): string {
		if (!/\$[^\$]+\$/.test(html)) {
			return html;
		};

		const reg = /(^|[^\d<\$]+)?\$((?:[^$<]|\.)*?)\$([^\d>\$]+|$)/gi;
		const tag = 'markuplatex';
		const code = 'markupcode';
		const regCode = new RegExp(`^${code}|${code}$`, 'i');
		const match = html.matchAll(reg);
		const render = (s: string) => {
			s = this.fromHtmlSpecialChars(s);

			let ret = s;
			try {
				ret = katex.renderToString(s, { 
					displayMode: false, 
					throwOnError: false,
					output: 'html',
					trust: ctx => [ '\\url', '\\href', '\\includegraphics' ].includes(ctx.command),
				});

				ret = ret ? ret : s;
			} catch (e) {};
			return ret;
		};

		let text = html;

		match.forEach((m: any) => {
			const m0 = String(m[0] || '');
			const m1 = String(m[1] || '');
			const m2 = String(m[2] || '');
			const m3 = String(m[3] || '');

			// Skip inline code marks
			if (regCode.test(m1) || regCode.test(m3)) {
				return;
			};

			// Skip Brazilian Real
			if (/R$/.test(m1) || /R$/.test(m2)) {
				return;
			};

			// Escaped $ sign
			if (/\\$/.test(m1) || /\\$/.test(m2)) {
				return;
			};

			text = text.replace(m0, `${m1}<${tag}>${render(m2)}</${tag}>${m3}`);
		});

		return text;
	};

};

export default new UtilCommon();