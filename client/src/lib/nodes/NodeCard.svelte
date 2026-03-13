<script lang="ts">
    import type { DashboardNode } from "../api/types";
    import NodeBars from "./NodeBars.svelte";
    import NodeContainers from "./NodeContainers.svelte";
    import NodeHeader from "./NodeHeader.svelte";

    type Props = {
        node: DashboardNode;
    };
    const { node }: Props = $props();
</script>

<div class="node-card">
    <NodeHeader {node} />
    <NodeBars {node} />
    <NodeContainers {node} />
</div>

<style>
    .node-card {
        background: var(--surface);
        border: 1px solid var(--border);
        border-radius: 10px;
        padding: 16px;
    }
</style>
