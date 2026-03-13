export function fmtBytes(bytes?: number): string {
    if (!bytes) return '—';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(0) + ' KB';
    if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(0) + ' MB';
    return (bytes / (1024 * 1024 * 1024)).toFixed(1) + ' GB';
}

export function fmtCPU(pct?: number): string {
    return pct ? pct.toFixed(1) + '%' : '—';
}

export function barClass(pct: number, base: string): string {
    if (pct >= 90) return base + ' crit';
    if (pct >= 70) return base + ' warn';
    return base;
}

export function shortImage(image: string): string {
    return image.replace(/@sha256:[0-9a-f]+/, '');
}

export async function fetchJSON<TData>(url: string) {
    const res = await fetch(url);
    if (!res.ok) throw new Error(`${url} → ${res.status}`);
    return res.json() as TData;
}
