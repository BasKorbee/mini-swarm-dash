<script lang="ts">
    import type { DashboardNode } from "../api/types";
    import NodeCard from "./NodeCard.svelte";

    type Props = {
        nodes: DashboardNode[];
    };
    const { nodes }: Props = $props();
</script>

<section>
    <h2>Nodes</h2>
    <div>
        {#each nodes as node (node.details.id ?? node.details.hostname)}
            <NodeCard {node} />
        {/each}
    </div>
</section>

<style>
    section div {
        display: flex;
        flex-direction: column;
        gap: 16px;
    }
</style>
