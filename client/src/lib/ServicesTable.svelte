<script lang="ts">
    import { shortImage } from './utils';

    export let services: any[] = [];
</script>

<table id="services-table">
    <thead>
        <tr>
            <th>Name</th>
            <th>Mode</th>
            <th>Replicas</th>
            <th>Image</th>
        </tr>
    </thead>
    <tbody>
        {#each services as s}
            {@const name     = s.Spec?.Name ?? '?'}
            {@const image    = shortImage(s.Spec?.TaskTemplate?.ContainerSpec?.Image ?? '?')}
            {@const isGlobal = !!s.Spec?.Mode?.Global}
            {@const replicas = isGlobal ? '—' : (s.Spec?.Mode?.Replicated?.Replicas ?? '?')}
            <tr>
                <td>{name}</td>
                <td>
                    {#if isGlobal}
                        <span class="tag global">global</span>
                    {:else}
                        <span class="tag replicated">replicated</span>
                    {/if}
                </td>
                <td>{replicas}</td>
                <td class="image-cell" title={image}>{image}</td>
            </tr>
        {/each}
    </tbody>
</table>
