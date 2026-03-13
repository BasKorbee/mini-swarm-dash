<script lang="ts">
    import type { ServiceSummary } from "../api/types";
    import ServicesTable from "./ServicesTable.svelte";

    type Props = {
        services: ServiceSummary[];
    };
    const { services }: Props = $props();
</script>

<section>
    <h2>Services</h2>
    <ServicesTable {services} />
</section>
