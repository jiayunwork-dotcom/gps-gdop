"use strict";

const $ = (id) => document.getElementById(id);

const rxX = $("rxX"), rxY = $("rxY"), rxZ = $("rxZ");
const maskInput = $("mask");
const satList = $("satList");
const errorBox = $("errorBox");

function showError(msg) {
  errorBox.textContent = msg;
  errorBox.hidden = false;
}

function clearError() {
  errorBox.hidden = true;
  errorBox.textContent = "";
}

function readRequest() {
  const sats = [];
  for (const line of satList.value.split("\n")) {
    const t = line.trim();
    if (t === "") continue;
    const parts = t.split(/[,\s]+/);
    if (parts.length !== 4) {
      throw new Error("卫星行需要 4 列：id, x, y, z — 得到「" + t + "」");
    }
    sats.push({ id: parts[0], x: Number(parts[1]), y: Number(parts[2]), z: Number(parts[3]) });
  }
  return {
    receiver_ecef: { x: Number(rxX.value), y: Number(rxY.value), z: Number(rxZ.value) },
    satellites: sats,
    elevation_mask_deg: Number(maskInput.value)
  };
}

async function postJSON(path, body) {
  const resp = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });
  const data = await resp.json();
  if (!resp.ok) {
    throw new Error(data.error || ("HTTP " + resp.status));
  }
  return data;
}

async function loadExample() {
  clearError();
  try {
    const examples = await (await fetch("api/examples")).json();
    const name = $("exampleSelect").value;
    const payload = examples[name];
    if (!payload) {
      throw new Error("示例 " + name + " 不存在");
    }
    rxX.value = payload.receiver_ecef.x;
    rxY.value = payload.receiver_ecef.y;
    rxZ.value = payload.receiver_ecef.z;
    maskInput.value = payload.elevation_mask_deg;
    satList.value = payload.satellites
      .map((s) => [s.id, s.x, s.y, s.z].join(", "))
      .join("\n");
  } catch (err) {
    showError("加载示例失败: " + err.message);
  }
}

function renderDop(d) {
  const tbody = document.querySelector("#dopTable tbody");
  tbody.innerHTML = "";
  const rows = [
    ["GDOP", d.gdop],
    ["PDOP", d.pdop],
    ["TDOP", d.tdop],
    ["HDOP", d.hdop],
    ["VDOP", d.vdop],
    ["条件数 cond(N)", d.condition_number],
    ["使用卫星数", d.satellites_used],
    ["被截止剔除", d.satellites_rejected]
  ];
  for (const [k, v] of rows) {
    const tr = document.createElement("tr");
    const th = document.createElement("td");
    th.textContent = k;
    const td = document.createElement("td");
    td.textContent = typeof v === "number" ? v.toFixed(4) : v;
    tr.append(th, td);
    tbody.appendChild(tr);
  }
  $("dopNote").textContent =
    "PDOP²+TDOP²=" + (d.pdop * d.pdop + d.tdop * d.tdop).toFixed(6) +
    "，GDOP²=" + (d.gdop * d.gdop).toFixed(6) +
    "，HDOP²+VDOP²=" + (d.hdop * d.hdop + d.vdop * d.vdop).toFixed(6);
  $("dopPanel").hidden = false;
}

async function computeDop() {
  clearError();
  let body;
  try {
    body = readRequest();
  } catch (err) {
    showError(err.message);
    return;
  }
  try {
    const d = await postJSON("api/dop", body);
    renderDop(d);
  } catch (err) {
    showError("计算 DOP 失败: " + err.message);
  }
}

function project(azDeg, elDeg) {
  // 天顶图投影：半径随仰角线性，90° 在圆心，0° 在圆周。
  const radius = (90 - elDeg) / 90 * 230;
  const az = azDeg * Math.PI / 180;
  return {
    x: radius * Math.sin(az),
    y: -radius * Math.cos(az)
  };
}

function renderSky(sky) {
  const svg = $("skySvg");
  svg.innerHTML = "";
  const NS = "http://www.w3.org/2000/svg";
  const make = (tag, attrs) => {
    const el = document.createElementNS(NS, tag);
    for (const k in attrs) el.setAttribute(k, attrs[k]);
    return el;
  };
  svg.appendChild(make("circle", { cx: 0, cy: 0, r: 230, fill: "#f4f7fa", stroke: "#333", "stroke-width": 1 }));
  for (const r of [76.7, 153.3]) {
    svg.appendChild(make("circle", { cx: 0, cy: 0, r: r, fill: "none", stroke: "#999", "stroke-dasharray": "3 3" }));
  }
  svg.appendChild(make("line", { x1: 0, y1: -230, x2: 0, y2: 230, stroke: "#999", "stroke-width": 0.5 }));
  svg.appendChild(make("line", { x1: -230, y1: 0, x2: 230, y2: 0, stroke: "#999", "stroke-width": 0.5 }));
  svg.appendChild(make("text", { x: 4, y: -236, "font-size": 12 })).textContent = "北 0°";
  svg.appendChild(make("text", { x: 236, y: 4, "font-size": 12 })).textContent = "东 90°";

  for (const p of sky.satellites) {
    const pt = project(p.azimuth_deg, p.elevation_deg);
    const g = make("g", {});
    g.appendChild(make("circle", {
      cx: pt.x, cy: pt.y, r: 7,
      fill: p.used ? "#1a73e8" : "#b0b0b0",
      stroke: "#000", "stroke-width": 0.5, opacity: p.used ? 1 : 0.5
    }));
    g.appendChild(make("text", {
      x: pt.x + 9, y: pt.y + 3, "font-size": 11
    })).textContent = p.id + (p.used ? "" : " (遮)");
    svg.appendChild(g);
  }
  $("skyNote").textContent = "仰角截止 " + sky.elevation_mask_deg + "°，使用 " +
    sky.used_count + " 颗，截止剔除 " + sky.rejected_count + " 颗。";
  $("skyPanel").hidden = false;
}

async function computeSky() {
  clearError();
  let body;
  try {
    body = readRequest();
  } catch (err) {
    showError(err.message);
    return;
  }
  try {
    const sky = await postJSON("api/sky", body);
    renderSky(sky);
  } catch (err) {
    showError("计算天顶图失败: " + err.message);
  }
}

$("loadExample").addEventListener("click", loadExample);
$("computeDop").addEventListener("click", computeDop);
$("computeSky").addEventListener("click", computeSky);
