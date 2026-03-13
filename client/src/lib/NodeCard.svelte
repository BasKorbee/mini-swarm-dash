<script lang="ts">
    import { fmtBytes, fmtCPU, barClass } from './utils';

    export let node: any;

    $: d = node.details;
    $: hostname  = d?.Description?.Hostname ?? '?';
    $: state     = d?.Status?.State ?? 'unknown';
    $: role      = d?.Spec?.Role ?? 'worker';
    $: isLeader  = d?.ManagerStatus?.Leader === true;
    $: nanoCPUs  = d?.Description?.Resources?.NanoCPUs ?? 0;
    $: memBytes  = d?.Description?.Resources?.MemoryBytes ?? 0;

    $: containers = node.containers ?? [];
    $: totalCPU = containers.reduce((s: number, c: any) => s + (c.cpu_percent ?? 0), 0);
    $: totalMem = containers.reduce((s: number, c: any) => s + (c.mem_usage ?? 0), 0);
    $: memPct   = memBytes > 0 ? (totalMem / memBytes) * 100 : 0;
    $: numCPUs  = nanoCPUs > 0 ? nanoCPUs / 1e9 : 1;
    $: cpuPct   = Math.min((totalCPU / (numCPUs * 100)) * 100, 100);

    $: dotClass  = state === 'ready' ? 'ready' : state === 'down' ? 'down' : 'drain';
    $: roleLabel = role === 'manager' ? (isLeader ? 'leader' : 'manager') : role;
    $: roleClass = isLeader ? 'leader' : role === 'manager' ? '' : 'worker';

    $: temp    = node.temp_celsius;
    $: tempPct = temp != null ? (temp / 90) * 100 : null;

    const containerCpuPct = (c: any) => c.cpu_percent ?? 0
    const containerMemPct = (c: any) => {
        const used = c.mem_usage ?? 0;
        return memBytes > 0 ? (used / memBytes) * 100 : 0;
    }
</script>

<div class="node-card">
    <div class="node-card-header">
        <span class="node-name">
            <span class="status-dot {dotClass}"></span>{hostname}
        </span>
        <span class="node-ip">{d?.Status?.Addr ?? ''}</span>
        <span class="node-role {roleClass}">{roleLabel}</span>
    </div>

    <div class="node-bars" style="--bar-cols:{tempPct != null ? 3 : 2}">
        <!-- CPU bar -->
        <div class="stat-row">
            <div class="stat-label">
                <span>CPU</span>
                <span>{fmtCPU(totalCPU)} / {numCPUs} cores</span>
            </div>
            <div class="bar-track">
                <div class="{barClass(cpuPct, 'bar-fill cpu')}" style="width:{Math.min(cpuPct, 100)}%"></div>
            </div>
        </div>

        <!-- Memory bar -->
        <div class="stat-row">
            <div class="stat-label">
                <span>Memory</span>
                <span>{fmtBytes(totalMem)} / {fmtBytes(memBytes)}</span>
            </div>
            <div class="bar-track">
                <div class="{barClass(memPct, 'bar-fill mem')}" style="width:{Math.min(memPct, 100)}%"></div>
            </div>
        </div>

        <!-- Temp bar (optional) -->
        {#if tempPct != null}
        <div class="stat-row">
            <div class="stat-label">
                <span>Temp</span>
                <span>{temp.toFixed(1)} °C</span>
            </div>
            <div class="bar-track">
                <div class="{barClass(tempPct, 'bar-fill temp')}" style="width:{Math.min(tempPct, 100)}%"></div>
            </div>
        </div>
        {/if}
    </div>

    <details class="node-containers" open>
        <summary class="node-containers-summary">
            {containers.length} container{containers.length !== 1 ? 's' : ''}
        </summary>
        <table class="containers-table">
            <thead>
                <tr><th>Service</th><th>Container</th><th>CPU %</th><th>Memory</th></tr>
            </thead>
            <tbody>
                {#if containers.length === 0}
                    <tr><td colspan="4" class="no-containers">No containers</td></tr>
                {:else}
                    {#each containers as c}
                        {@const cpuP = containerCpuPct(c)}
                        {@const memP = containerMemPct(c)}
                        {@const service = c.service_name ?? c.name?.split('.')[0] ?? '?'}
                        <tr>
                            <td>{service}</td>
                            <td style="color:var(--muted);font-size:12px">{c.name ?? '?'}</td>
                            <td>
                                <div class="inline-bar">
                                    <div class="bar-track">
                                        <div class="{barClass(cpuP, 'bar-fill cpu')}" style="width:{Math.min(cpuP, 100)}%"></div>
                                    </div>
                                    <span>{fmtCPU(cpuP)}</span>
                                </div>
                            </td>
                            <td>
                                <div class="inline-bar">
                                    <div class="bar-track">
                                        <div class="{barClass(memP, 'bar-fill mem')}" style="width:{Math.min(memP, 100)}%"></div>
                                    </div>
                                    <span>{fmtBytes(c.mem_usage)}</span>
                                </div>
                            </td>
                        </tr>
                    {/each}
                {/if}
            </tbody>
        </table>
    </details>
</div>
