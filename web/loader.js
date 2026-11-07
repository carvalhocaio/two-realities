const MB = 1024 * 1024;

const TEXT = {
  en: { loading: "LOADING", starting: "STARTING", failed: "FAILED TO LOAD" },
  pt: { loading: "CARREGANDO", starting: "PREPARANDO", failed: "FALHA AO CARREGAR" },
};

export function language(navigatorLanguage) {
  return String(navigatorLanguage || "").toLowerCase().startsWith("pt") ? "pt" : "en";
}

export function expectedSize(headers, injected) {
  const size = Number(injected);
  if (Number.isFinite(size) && size > 0) {
    return size;
  }
  const encoding = (headers.get("content-encoding") || "identity").toLowerCase();
  if (encoding !== "identity") {
    return 0;
  }
  const length = Number(headers.get("content-length"));
  return Number.isFinite(length) && length > 0 ? length : 0;
}

export function fraction(received, total) {
  if (!(total > 0)) {
    return null;
  }
  return Math.min(Math.max(received / total, 0), 1);
}

export function label(lang, stage, received, total) {
  const t = TEXT[lang] || TEXT.en;
  if (stage === "starting") {
    return `${t.starting}…`;
  }
  if (stage === "failed") {
    return t.failed;
  }
  const f = fraction(received, total);
  if (f === null && !(received > 0)) {
    return `${t.loading}…`;
  }
  if (f === null) {
    return `${t.loading} ${(received / MB).toFixed(1)} MB`;
  }
  return `${t.loading} ${Math.floor(f * 100)}%`;
}

export async function download(response, total, onProgress) {
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}`);
  }
  if (!response.body || typeof response.body.getReader !== "function") {
    const bytes = new Uint8Array(await response.arrayBuffer());
    onProgress(bytes.length, total);
    return bytes;
  }
  const reader = response.body.getReader();
  const chunks = [];
  let received = 0;
  for (;;) {
    const { done, value } = await reader.read();
    if (done) {
      break;
    }
    chunks.push(value);
    received += value.length;
    onProgress(received, total);
  }
  const bytes = new Uint8Array(received);
  let offset = 0;
  for (const chunk of chunks) {
    bytes.set(chunk, offset);
    offset += chunk.length;
  }
  return bytes;
}
