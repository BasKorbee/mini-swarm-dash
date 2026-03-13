<script lang="ts">
    import type { DashboardNode } from "../api/types";
    import NodeContainer from "./NodeContainer.svelte";

    type Props = {
        node: DashboardNode;
    };
    const { node }: Props = $props();

    const containers = $derived(node.containers ?? []);
</script>

<details open>
    <summary>
        {containers.length} container{containers.length !== 1 ? "s" : ""}
    </summary>
    <table>
        <colgroup>
            <col style="width:28%" />
            <col style="width:30%" />
            <col style="width:21%" />
            <col style="width:21%" />
        </colgroup>
        <thead>
            <tr
                ><th>Service</th><th>Container</th><th>CPU %</th><th>Memory</th
                ></tr
            >
        </thead>
        <tbody>
            {#if containers.length === 0}
                <tr><td colspan="4" class="no-containers">No containers</td></tr
                >
            {:else}
                {#each containers as container}
                    <NodeContainer {node} {container} />
                {/each}
            {/if}
        </tbody>
    </table>
</details>

<style>
    details {
        margin-top: 14px;
        border-top: 1px solid var(--border);
        padding-top: 10px;
    }

    details summary {
        font-size: 12px;
        color: var(--muted);
        cursor: pointer;
        user-select: none;
        list-style: none;
        display: flex;
        align-items: center;
        gap: 6px;
        margin-bottom: 0;
    }

    details summary::-webkit-details-marker {
        display: none;
    }

    details summary::before {
        content: "▶";
        font-size: 9px;
        color: var(--muted);
        transition: transform 0.2s ease;
    }

    details[open] summary::before {
        transform: rotate(90deg);
    }

    details[open] summary {
        margin-bottom: 10px;
    }

    /* Inner containers table */
    details table {
        width: 100%;
        table-layout: fixed;
        border-collapse: collapse;
        background: transparent;
        border: none;
        border-radius: 0;
    }

    details table thead th {
        background: transparent;
        color: var(--muted);
        font-size: 10px;
        text-transform: uppercase;
        letter-spacing: 0.8px;
        padding: 4px 8px;
        text-align: left;
        border-bottom: 1px solid var(--border);
    }

    details table tbody td {
        padding: 6px 8px;
        border-bottom: 1px solid
            color-mix(in srgb, var(--border) 50%, transparent);
        color: var(--text);
        font-size: 12px;
    }

    details table tbody tr:last-child td {
        border-bottom: none;
    }
    details table tbody tr:hover {
        background: rgba(255, 255, 255, 0.03);
    }

    .no-containers {
        color: var(--muted);
        font-size: 12px;
        font-style: italic;
        padding: 6px 8px;
    }
</style>
