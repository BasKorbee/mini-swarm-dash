<script lang="ts">
    import { barClass } from "./utils";

    type Props = {
        title?: string;
        numerator: string;
        denominator?: string;
        percent: number;
        base: "accent" | "green" | "orange" | "warn" | "crit";
        inline?: boolean;
    };
    const { title, numerator, denominator, percent, base, inline }: Props =
        $props();
    const fraction = $derived(
        `${numerator}${denominator ? ` / ${denominator}` : ""}`,
    );

    const barClassString = $derived(barClass(percent, `bar-fill ${base}`));
    const width = $derived(`width:${Math.min(percent, 100)}%`);
</script>

{#if inline}
    <div class="inline-bar">
        <div class="bar-track">
            <div class={barClassString} style={width}></div>
        </div>
        <span>{numerator}</span>
    </div>
{:else}
    <div class="stat-row">
        <div class="stat-label">
            {#if title}
                <span>{title}</span>
            {/if}
            <span>{fraction}</span>
        </div>
        <div class="bar-track">
            <div class={barClassString} style={width}></div>
        </div>
    </div>
{/if}

<style>
    .stat-row {
        margin-bottom: 10px;
    }

    .stat-label {
        display: flex;
        justify-content: space-between;
        font-size: 12px;
        color: var(--muted);
        margin-bottom: 4px;
    }

    .stat-label span:last-child {
        color: var(--text);
        font-variant-numeric: tabular-nums;
    }

    .bar-track {
        height: 6px;
        background: var(--bar-track);
        border-radius: 3px;
        overflow: hidden;
    }

    .bar-fill {
        height: 100%;
        border-radius: 3px;
        transition: width 0.4s ease;
    }

    .bar-fill.accent {
        background: var(--accent);
    }
    .bar-fill.green {
        background: var(--green);
    }
    .bar-fill.orange {
        background: var(--orange);
    }
    .bar-fill.warn {
        background: var(--yellow);
    }
    .bar-fill.crit {
        background: var(--red);
    }

    .inline-bar {
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .inline-bar .bar-track {
        flex: 1;
        max-width: 80px;
    }

    .inline-bar span {
        font-variant-numeric: tabular-nums;
        min-width: 40px;
        text-align: right;
        font-size: 12px;
    }
</style>
