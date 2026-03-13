<script lang="ts">
    import type { DashboardNode } from "../api/types";
    import NodeRole from "./NodeRole.svelte";
    import NodeTitle from "./NodeTitle.svelte";

    type Props = {
        node: DashboardNode
    }
    const { node }: Props = $props()
</script>

<div class="node-card-header">
    <NodeTitle {node} />
    <NodeRole {node} />
</div>

<style>
    .node-card-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 4px;
        margin-bottom: 14px;
    }
</style>
