<script lang="ts">
    import Logo from "./Logo.svelte";
    import Pause from "./Pause.svelte";
    import Play from "./Play.svelte";

    type Props = {
        lastUpdated: string;
        refreshRate: number;
        onRateChange: () => void;
        paused?: boolean;
        togglePause: () => void;
    };
    let { lastUpdated, refreshRate, onRateChange, paused, togglePause }: Props =
        $props();
</script>

<header>
    <h1><Logo width={16} height={16} /> Mini Swarm Dashboard</h1>
    <div>
        <span>Updated {lastUpdated} </span>
        <select bind:value={refreshRate} onchange={onRateChange}>
            <option value={1}>1s</option>
            <option value={5}>5s</option>
            <option value={10}>10s</option>
            <option value={30}>30s</option>
            <option value={60}>1m</option>
            <option value={300}>5m</option>
        </select>
        <button
            title={paused ? "Resume refreshing" : "Pause refreshing"}
            onclick={togglePause}
        >
            {#if paused}
                <Play />
            {:else}
                <Pause />
            {/if}
        </button>
    </div>
</header>

<style>
    header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 24px;
        border-bottom: 1px solid var(--border);
        background: var(--surface);
    }

    header h1 {
        font-size: 18px;
        font-weight: 600;
        letter-spacing: 0.5px;
        display: flex;
        align-items: center;
        gap: 4px;
    }

    header div {
        display: flex;
        justify-content: end;
        align-items: center;
        flex-wrap: wrap;
        gap: 4px;
    }

    header span {
        font-size: 12px;
        color: var(--muted);
    }

    header select {
        background: none;
        border: 1px solid var(--border);
        color: var(--muted);
        border-radius: 6px;
        padding: 2px 6px;
        font-size: 14px;
        cursor: pointer;
        line-height: 1.6;
        height: 30px;
        box-sizing: border-box;
    }

    header select:hover {
        color: var(--text);
        border-color: var(--text);
    }

    header button {
        background: none;
        border: 1px solid var(--border);
        color: var(--muted);
        border-radius: 6px;
        padding: 2px 8px;
        font-size: 14px;
        cursor: pointer;
        line-height: 1.6;
        height: 30px;
        box-sizing: border-box;
    }

    header button:hover {
        color: var(--text);
        border-color: var(--text);
    }
</style>
