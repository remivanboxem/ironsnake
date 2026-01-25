<script lang="ts">
	import { onMount } from 'svelte';
	import MarkdownIt from 'markdown-it';
	import katex from 'katex';

	interface Props {
		content: string;
		class?: string;
	}

	let { content, class: className = '' }: Props = $props();

	let renderedHtml = $state('');

	// Custom renderer for math expressions
	function renderMath(tex: string, displayMode: boolean): string {
		try {
			return katex.renderToString(tex, {
				displayMode,
				throwOnError: false,
				strict: false
			});
		} catch (e) {
			console.error('KaTeX error:', e);
			return `<span class="text-red-500">${tex}</span>`;
		}
	}

	// Process :math:`...` syntax and convert to rendered math
	function processMathSyntax(text: string): string {
		// Handle :math:`...` inline math (RST/Sphinx style)
		text = text.replace(/:math:`([^`]+)`/g, (_, tex) => {
			return renderMath(tex, false);
		});

		// Handle $$...$$ display math
		text = text.replace(/\$\$([^$]+)\$\$/g, (_, tex) => {
			return `<div class="my-2">${renderMath(tex.trim(), true)}</div>`;
		});

		// Handle $...$ inline math (but not escaped \$)
		text = text.replace(/(?<!\\)\$([^$\n]+)\$/g, (_, tex) => {
			return renderMath(tex, false);
		});

		return text;
	}

	// Initialize markdown-it
	const md = new MarkdownIt({
		html: false,
		linkify: true,
		typographer: true,
		breaks: true
	});

	// Custom code block renderer for syntax highlighting styling
	const defaultFence = md.renderer.rules.fence;
	md.renderer.rules.fence = (tokens, idx, options, env, self) => {
		const token = tokens[idx];
		const lang = token.info.trim();
		const code = token.content;

		return `<pre class="bg-muted p-4 rounded-md overflow-x-auto my-2"><code class="language-${lang || 'text'} text-sm">${md.utils.escapeHtml(code)}</code></pre>`;
	};

	// Custom inline code renderer
	const defaultCodeInline = md.renderer.rules.code_inline;
	md.renderer.rules.code_inline = (tokens, idx, options, env, self) => {
		const token = tokens[idx];
		return `<code class="bg-muted px-1.5 py-0.5 rounded text-sm font-mono">${md.utils.escapeHtml(token.content)}</code>`;
	};

	$effect(() => {
		if (content) {
			// First process math syntax, then render markdown
			const processedContent = processMathSyntax(content);
			// For content that already has math rendered (HTML), we need to be careful
			// Render markdown on the original, then process math on the result
			let html = md.render(content);
			// Now process math in the rendered HTML
			html = processMathSyntax(html);
			renderedHtml = html;
		} else {
			renderedHtml = '';
		}
	});
</script>

<svelte:head>
	<link
		rel="stylesheet"
		href="https://cdn.jsdelivr.net/npm/katex@0.16.27/dist/katex.min.css"
		integrity="sha384-AKaEXjCwKFX7zegKe5EK7hDm2xdlS0zPgEfvKa0r7CMd/jVk6x7xLoNTCcPFEELw"
		crossorigin="anonymous"
	/>
</svelte:head>

<div class="markdown-content prose prose-sm max-w-none dark:prose-invert {className}">
	{@html renderedHtml}
</div>

<style>
	.markdown-content :global(p) {
		margin-bottom: 0.75rem;
	}

	.markdown-content :global(p:last-child) {
		margin-bottom: 0;
	}

	.markdown-content :global(ul),
	.markdown-content :global(ol) {
		margin-left: 1.5rem;
		margin-bottom: 0.75rem;
	}

	.markdown-content :global(li) {
		margin-bottom: 0.25rem;
	}

	.markdown-content :global(h1),
	.markdown-content :global(h2),
	.markdown-content :global(h3) {
		font-weight: 600;
		margin-top: 1rem;
		margin-bottom: 0.5rem;
	}

	.markdown-content :global(blockquote) {
		border-left: 3px solid hsl(var(--muted));
		padding-left: 1rem;
		margin: 0.75rem 0;
		color: hsl(var(--muted-foreground));
	}

	.markdown-content :global(a) {
		color: hsl(var(--primary));
		text-decoration: underline;
	}

	.markdown-content :global(a:hover) {
		text-decoration: none;
	}

	/* KaTeX specific styles */
	.markdown-content :global(.katex) {
		font-size: 1.1em;
	}

	.markdown-content :global(.katex-display) {
		margin: 0.75rem 0;
		overflow-x: auto;
	}
</style>
