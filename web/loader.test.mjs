import { test } from "node:test";
import assert from "node:assert/strict";

import { download, expectedSize, fraction, label, language } from "./loader.js";

const headers = (h) => new Headers(h);

test("language", () => {
  for (const [input, want] of [
    ["pt-BR", "pt"],
    ["PT", "pt"],
    ["en-US", "en"],
    ["", "en"],
    [undefined, "en"],
  ]) {
    assert.equal(language(input), want, String(input));
  }
});

test("expectedSize", () => {
  for (const [name, h, injected, want] of [
    ["injected size wins over headers", { "content-length": "10" }, "500", 500],
    ["injected size wins over gzip", { "content-encoding": "gzip", "content-length": "10" }, 500, 500],
    ["zero injected falls back to content-length", { "content-length": "42" }, "0", 42],
    ["garbage injected falls back to content-length", { "content-length": "42" }, "abc", 42],
    ["identity encoding keeps content-length", { "content-encoding": "identity", "content-length": "42" }, undefined, 42],
    ["gzip length is the compressed size, so unknown", { "content-encoding": "gzip", "content-length": "42" }, undefined, 0],
    ["missing content-length is unknown", {}, undefined, 0],
    ["invalid content-length is unknown", { "content-length": "x" }, undefined, 0],
  ]) {
    assert.equal(expectedSize(headers(h), injected), want, name);
  }
});

test("fraction", () => {
  for (const [received, total, want] of [
    [50, 200, 0.25],
    [0, 200, 0],
    [300, 200, 1],
    [-5, 200, 0],
    [10, 0, null],
    [10, undefined, null],
  ]) {
    assert.equal(fraction(received, total), want, `${received}/${total}`);
  }
});

test("label", () => {
  for (const [args, want] of [
    [["en", "loading", 42, 100], "LOADING 42%"],
    [["pt", "loading", 999, 1000], "CARREGANDO 99%"],
    [["en", "loading", 1000, 1000], "LOADING 100%"],
    [["en", "loading", 3 * 1024 * 1024, 0], "LOADING 3.0 MB"],
    [["pt", "starting", 0, 0], "PREPARANDO…"],
    [["en", "starting", 0, 0], "STARTING…"],
    [["pt", "failed", 0, 0], "FALHA AO CARREGAR"],
    [["xx", "loading", 1, 2], "LOADING 50%"],
    [["pt", "loading", 0, 0], "CARREGANDO…"],
    [["en", "loading", 0, 100], "LOADING 0%"],
  ]) {
    assert.equal(label(...args), want, args.join(" "));
  }
});

const streamOf = (chunks) =>
  new ReadableStream({
    start(c) {
      for (const chunk of chunks) {
        c.enqueue(new Uint8Array(chunk));
      }
      c.close();
    },
  });

test("download streams the body and reports progress", async () => {
  const calls = [];
  const bytes = await download(new Response(streamOf([[1, 2], [3], [4, 5, 6]])), 6, (r, t) => calls.push([r, t]));
  assert.deepEqual([...bytes], [1, 2, 3, 4, 5, 6]);
  assert.deepEqual(calls, [[2, 6], [3, 6], [6, 6]]);
});

test("download of an empty body", async () => {
  const calls = [];
  const bytes = await download(new Response(streamOf([])), 0, (r) => calls.push(r));
  assert.equal(bytes.length, 0);
  assert.deepEqual(calls, []);
});

test("download without a readable stream falls back to arrayBuffer", async () => {
  const calls = [];
  const fake = { ok: true, body: null, arrayBuffer: async () => new Uint8Array([7, 8]).buffer };
  const bytes = await download(fake, 2, (r, t) => calls.push([r, t]));
  assert.deepEqual([...bytes], [7, 8]);
  assert.deepEqual(calls, [[2, 2]]);
});

test("download rejects an HTTP error", async () => {
  await assert.rejects(download(new Response("nope", { status: 404 }), 0, () => {}), /HTTP 404/);
});

test("download rejects a broken stream", async () => {
  const broken = new ReadableStream({
    pull(c) {
      c.error(new Error("connection reset"));
    },
  });
  await assert.rejects(download(new Response(broken), 10, () => {}), /connection reset/);
});
