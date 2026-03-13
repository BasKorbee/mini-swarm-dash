<script lang="ts">
    import type { ServiceSummary } from '../api/types';
    import ServicesTableRow from './ServicesTableRow.svelte';

    type Props = {
        services: ServiceSummary[]
    }
    const { services }: Props = $props()
</script>

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Mode</th>
            <th>Replicas</th>
            <th>Image</th>
        </tr>
    </thead>
    <tbody>
        {#each services as service}
            <ServicesTableRow {service} />
        {/each}
    </tbody>
</table>

<style>

</style>