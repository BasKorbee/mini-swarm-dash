<script lang="ts">
    let { width = 12, height = 12 } = $props() 
</script>

<svg
    xmlns="http://www.w3.org/2000/svg"
    {width}
    {height}
    viewBox="0 0 24 24"
    ><path
        fill="none"
        stroke="currentColor"
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M5 4h4v16H5Zm10 0h4v16h-4Z"
    /></svg
>
