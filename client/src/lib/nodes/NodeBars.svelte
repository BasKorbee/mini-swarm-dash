<script lang="ts">
    import type { DashboardNode } from "../api/types";
    import Bar from "../Bar.svelte";
    import { fmtBytes, fmtCPU } from "../utils";

    type Props = {
        node: DashboardNode;
    };
    const { node }: Props = $props();

    const containers = $derived(node.containers ?? []);
    const totalCPU = $derived(
        containers.reduce((s, c) => s + (c.cpu_percent ?? 0), 0),
    );
    const totalMemory = $derived(
        containers.reduce((s, c) => s + (c.mem_usage ?? 0), 0),
    );
    const memoryPercent = $derived(
        node.details.memory_bytes > 0
            ? (totalMemory / node.details.memory_bytes) * 100
            : 0,
    );
    const numberOfCPUs = $derived(
        (node.details.nano_cpus ?? 0) > 0
            ? (node.details.nano_cpus ?? 0) / 1e9
            : 1,
    );
    const cpuPercent = $derived(
        Math.min((totalCPU / (numberOfCPUs * 100)) * 100, 100),
    );

    const tempPercent = $derived(
        node.temp_celsius != null ? (node.temp_celsius / 90) * 100 : null,
    );
</script>

<div class="node-bars" style="--bar-cols:{tempPercent != null ? 3 : 2}">
    <Bar
        title="CPU"
        numerator={fmtCPU(totalCPU)}
        denominator={numberOfCPUs + " cores"}
        percent={cpuPercent}
        base="accent"
    />
    <Bar
        title="Memory"
        numerator={fmtBytes(totalMemory)}
        denominator={fmtBytes(node.details.memory_bytes ?? 0)}
        percent={memoryPercent}
        base="green"
    />
    {#if tempPercent != null && node.temp_celsius != null}
        <Bar
            title="Temp"
            numerator={node.temp_celsius.toFixed(1) + "°C"}
            percent={tempPercent}
            base="orange"
        />
    {/if}
</div>

<style>
    .node-bars {
        display: grid;
        grid-template-columns: 1fr;
        gap: 0 24px;
    }

    @media (min-width: 900px) {
        .node-bars {
            grid-template-columns: repeat(var(--bar-cols, 2), 1fr);
        }
    }
</style>
