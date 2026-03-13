<script lang="ts">
    import type { ServiceSummary } from "../api/types";
    import { shortImage } from "../utils";

    type Props = {
        service: ServiceSummary;
    };
    const { service }: Props = $props();

    const name = $derived(service.name ?? "?");
    const image = $derived(shortImage(service.image ?? "?"));
</script>

<tr>
    <td>{name}</td>
    <td>
        {#if service.is_global}
            <span class="tag global">global</span>
        {:else}
            <span class="tag replicated">replicated</span>
        {/if}
    </td>
    <td>{service.is_global ? "-" : (service.replicas ?? "?")}</td>
    <td class="image-cell" title={image}>{image}</td>
</tr>

<style>
    .image-cell {
        max-width: 300px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        color: var(--muted);
        font-size: 12px;
        font-family: monospace;
    }

    .tag {
        display: inline-block;
        padding: 1px 7px;
        border-radius: 99px;
        font-size: 11px;
    }

    .tag.global {
        background: #2d1b4e;
        color: #c084fc;
    }
    .tag.replicated {
        background: #1e3a5f;
        color: var(--accent);
    }
</style>
