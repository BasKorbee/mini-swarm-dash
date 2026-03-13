<script lang="ts">
    import type { DashboardNode } from "../api/types";

    type Props = {
        node: DashboardNode
    };
    const { node }: Props = $props();
</script>

<div>
    <span class="node-name">
        <span class="status-dot {node.details.state ?? 'drain'}"></span>{node.details.hostname ?? "?"}
    </span>
    <span class="node-ip">{node.details.addr ?? ""}</span>
</div>

<style>
    .node-name {
        font-weight: 600;
        font-size: 15px;
        min-width: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
    .status-dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        display: inline-block;
        margin-right: 6px;
    }
    .status-dot.ready {
        background: var(--green);
    }
    .status-dot.down {
        background: var(--red);
    }
    .status-dot.drain {
        background: var(--yellow);
    }
    .node-ip {
        font-size: 11px;
        color: var(--muted);
        flex: 1;
    }
</style>
