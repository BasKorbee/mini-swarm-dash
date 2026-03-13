<script lang="ts">
    import type { DashboardNode } from "../api/types";

    type Props = {
        node: DashboardNode;
    };
    const { node }: Props = $props();

    function getRoleLabel(role: string | undefined, leader: boolean): string {
        if (leader) return "leader";
        if (role === "manager") return "manager";
        return role ?? "worker";
    }
    function getRoleClass(role: string | undefined, leader: boolean): string {
        if (leader) return "leader";
        if (role === "manager") return "manager";
        return "worker";
    }
    const roleLabel = $derived(getRoleLabel(node.details.role, node.details.is_leader));
    const roleClass = $derived(getRoleClass(node.details.role, node.details.is_leader));
</script>

<span class="node-role {roleClass}">{roleLabel}</span>

<style>
    .node-role {
        font-size: 11px;
        padding: 2px 7px;
        border-radius: 99px;
        background: #1e3a5f;
        color: var(--accent);
    }

    .node-role.worker {
        background: #1e2d1e;
        color: var(--green);
    }

    .node-role.leader {
        background: #2d2000;
        color: #fbbf24;
    }
</style>
