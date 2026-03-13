<script lang="ts">
	import { onMount, onDestroy } from "svelte";
	import type { DashboardNode, ServiceSummary } from './lib/api/types'
	import NodeCard from "./lib/nodes/NodeCard.svelte";
	import ServicesTable from "./lib/services/ServicesTable.svelte";
	import { fetchJSON } from "./lib/utils";
    import Header from "./lib/Header.svelte";
    import Nodes from "./lib/nodes/Nodes.svelte";
    import Services from "./lib/services/Services.svelte";

	let nodes: DashboardNode[] = [];
	let services: ServiceSummary[] = [];
	let lastUpdated = "";
	let paused = false;
	let intervalId: ReturnType<typeof setInterval>;
	let refreshRate = 5;

	async function refresh() {
		try {
			[nodes, services] = await Promise.all([
				fetchJSON<DashboardNode[]>("/api/nodes"),
				fetchJSON<ServiceSummary[]>("/api/services"),
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
	<Nodes {nodes} />
	<Services {services} />
</main>
