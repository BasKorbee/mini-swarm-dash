<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import NodeCard from "./lib/NodeCard.svelte";
	import ServicesTable from "./lib/ServicesTable.svelte";
	import { fetchJSON } from "./lib/utils";
    import Header from "./lib/Header.svelte";

	let nodes: any[] = [];
	let services: any[] = [];
	let lastUpdated = "";
	let paused = false;
	let intervalId: ReturnType<typeof setInterval>;
	let refreshRate = 5;

	async function refresh() {
		try {
			[nodes, services] = await Promise.all([
				fetchJSON("/api/nodes"),
				fetchJSON("/api/services"),
			]);
			lastUpdated = new Date().toLocaleTimeString();
		} catch (err) {
			console.error("Dashboard refresh error:", err);
		}
	}

	function togglePause() {
		paused = !paused;
		if (paused) {
			clearInterval(intervalId);
		} else {
			refresh();
			intervalId = setInterval(refresh, refreshRate * 1000);
		}
	}

	function onRateChange() {
		if (!paused) {
			clearInterval(intervalId);
			intervalId = setInterval(refresh, refreshRate * 1000);
		}
	}

	onMount(() => {
		refresh();
		intervalId = setInterval(refresh, refreshRate * 1000);
	});

	onDestroy(() => clearInterval(intervalId));
</script>

<Header {lastUpdated} {refreshRate} {onRateChange} {paused} {togglePause} />

<main>
	<section>
		<h2>Nodes</h2>
		<div id="nodes-grid">
			{#each nodes as node (node.details?.ID ?? node.details?.Description?.Hostname)}
				<NodeCard {node} />
			{/each}
		</div>
	</section>

	<section>
		<h2>Services</h2>
		<ServicesTable {services} />
	</section>
</main>
