<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import NodeCard from "./lib/NodeCard.svelte";
	import ServicesTable from "./lib/ServicesTable.svelte";
	import { fetchJSON } from "./lib/utils";

	let nodes: any[] = [];
	let services: any[] = [];
	let lastUpdated = "";
	let paused = false;
	let intervalId: ReturnType<typeof setInterval>;

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
			intervalId = setInterval(refresh, 5000);
		}
	}

	onMount(() => {
		refresh();
		intervalId = setInterval(refresh, 5000);
	});

	onDestroy(() => clearInterval(intervalId));
</script>

<header>
	<h1>&#9783; Mini Swarm Dashboard</h1>
	<div>
		<span id="last-updated">Updated {lastUpdated} </span>
		<button
			id="pause-btn"
			title={paused ? "Resume refreshing" : "Pause refreshing"}
			on:click={togglePause}>{paused ? "▶" : "⏸"}</button
		>
	</div>
</header>

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
