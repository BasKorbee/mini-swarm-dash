<script lang="ts">
    import type { ContainerStats, DashboardNode } from "../api/types";
    import Bar from "../Bar.svelte";
    import { fmtBytes, fmtCPU } from "../utils";

    type Props = {
        node: DashboardNode;
        container: ContainerStats;
    };
    const { node, container }: Props = $props();

    const containerCpuPercent = $derived(container.cpu_percent ?? 0);
    const containerMemoryPercent = $derived(
        node.details.memory_bytes > 0
            ? (container.mem_usage / node.details.memory_bytes) * 100
            : 0,
    );

    const service = $derived(
        container.service_name ?? container.name?.split(".")[0] ?? "?",
    );
    const name = $derived(container.name ?? "?");
</script>

<tr>
    <td class="service" title={service}>{service}</td>
    <td class="name" title={name}>{name}</td>
    <td>
        <Bar
            numerator={fmtCPU(containerCpuPercent)}
            percent={containerCpuPercent}
            base="accent"
            inline
        />
    </td>
    <td>
        <Bar
            numerator={fmtBytes(container.mem_usage)}
            percent={containerMemoryPercent}
            base="green"
            inline
        />
    </td>
</tr>

<style>
    .service {
        max-width: 100px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .name {
        color: var(--muted);
        font-size: 12px;
        max-width: 120px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
</style>
