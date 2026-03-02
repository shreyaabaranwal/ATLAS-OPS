'use client'





import { useState, useEffect, useRef } from "react";

import { motion, useInView } from "framer-motion";
import {
  AreaChart, Area, LineChart, Line, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, RadialBarChart, RadialBar
} from "recharts";

// ─── THEME ────────────────────────────────────────────────────────────────────
const ORANGE = "#FF6B2B";
const ORANGE_DIM = "#cc4d1a";
const ORANGE_GLOW = "rgba(255,107,43,0.4)";
const WHITE = "#FFFFFF";
const DARK_BG = "#070A0F";
const CARD_BG = "rgba(255,255,255,0.04)";
const CARD_BORDER = "rgba(255,107,43,0.2)";

// ─── MOCK DATA ─────────────────────────────────────────────────────────────────
const cpuHistory = Array.from({ length: 60 }, (_, i) => {
  const base = 45 + Math.sin(i * 0.3) * 20;
  return { time: `${60 - i}s`, cpu: Math.round(base + Math.random() * 10), ema: Math.round(base * 0.9 + Math.random() * 8) };
}).reverse();

const incidents = [
  { id: "INC-0041", state: "PENDING",   attempts: 1, confidence: 87 },
  { id: "INC-0040", state: "RESOLVED",  attempts: 3, confidence: 94 },
  { id: "INC-0039", state: "FAILED",    attempts: 5, confidence: 31 },
  { id: "INC-0038", state: "EXECUTING", attempts: 2, confidence: 78 },
  { id: "INC-0037", state: "RESOLVED",  attempts: 1, confidence: 92 },
];

const timeline = [
  { action: "Scale Kubernetes Deployment", time: "14:32:01", cpuBefore: 89, cpuAfter: 52, status: "success" },
  { action: "Restart Overloaded Pod",      time: "14:28:44", cpuBefore: 76, cpuAfter: 41, status: "success" },
  { action: "Throttle Ingress Traffic",    time: "14:21:17", cpuBefore: 91, cpuAfter: 88, status: "partial" },
  { action: "Flush Redis Cache",           time: "14:15:08", cpuBefore: 73, cpuAfter: 38, status: "success" },
  { action: "Auto-rollback Deployment",    time: "14:09:33", cpuBefore: 95, cpuAfter: 61, status: "rollback" },
];

const dlqItems = [
  { id: "DLQ-007", error: "Timeout: k8s API unreachable",   ts: "14:30:11" },
  { id: "DLQ-006", error: "Auth failure: IAM role expired", ts: "13:58:22" },
];

const metricCards = [
  { label: "Total Incidents", value: 247,  suffix: "",  color: ORANGE },
  { label: "Success Rate",    value: 94.2, suffix: "%", color: "#4AFFC4" },
  { label: "Rollbacks",       value: 12,   suffix: "",  color: "#FF4A6B" },
  { label: "Avg Exec Time",   value: 3.4,  suffix: "s", color: "#9B8BFF" },
];

const navItems = ["Monitor","Analytics","Incidents","Timeline","DLQ","Metrics"];

const stateConfig = {
  PENDING:   { color: "#FFB347", glow: "rgba(255,179,71,0.4)" },
  RESOLVED:  { color: "#4AFFC4", glow: "rgba(74,255,196,0.4)" },
  FAILED:    { color: "#FF4A6B", glow: "rgba(255,74,107,0.4)" },
  EXECUTING: { color: ORANGE,    glow: ORANGE_GLOW },
};

const statusColors = { success: "#4AFFC4", partial: "#FFB347", rollback: "#FF4A6B" };

// ─── PARTICLE CANVAS ──────────────────────────────────────────────────────────
function ParticleCanvas() {
  const canvasRef = useRef<HTMLCanvasElement | null>(null)

  useEffect(() => {
    const canvas = canvasRef.current
    if (!canvas) return

    const ctx = canvas.getContext("2d")
    if (!ctx) return

    let W = (canvas.width = window.innerWidth)
    let H = (canvas.height = window.innerHeight)

    const pts = Array.from({ length: 100 }, () => ({
      x: Math.random() * W,
      y: Math.random() * H,
      r: Math.random() * 1.4 + 0.3,
      vx: (Math.random() - 0.5) * 0.25,
      vy: (Math.random() - 0.5) * 0.25,
    }))

    let raf: number

    const draw = () => {
      ctx.clearRect(0, 0, W, H)

      pts.forEach(p => {
        p.x += p.vx
        p.y += p.vy

        if (p.x < 0) p.x = W
        if (p.x > W) p.x = 0
        if (p.y < 0) p.y = H
        if (p.y > H) p.y = 0

        ctx.beginPath()
        ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2)
        ctx.fillStyle = "rgba(255,107,43,0.35)"
        ctx.fill()
      })

      for (let i = 0; i < pts.length; i++) {
        for (let j = i + 1; j < pts.length; j++) {
          const dx = pts[i].x - pts[j].x
          const dy = pts[i].y - pts[j].y
          const d = Math.sqrt(dx * dx + dy * dy)

          if (d < 110) {
            ctx.beginPath()
            ctx.moveTo(pts[i].x, pts[i].y)
            ctx.lineTo(pts[j].x, pts[j].y)
            ctx.strokeStyle = `rgba(255,107,43,${0.07 * (1 - d / 110)})`
            ctx.stroke()
          }
        }
      }

      raf = requestAnimationFrame(draw)
    }

    draw()

    const onResize = () => {
      W = canvas.width = window.innerWidth
      H = canvas.height = window.innerHeight
    }

    window.addEventListener("resize", onResize)

    return () => {
      cancelAnimationFrame(raf)
      window.removeEventListener("resize", onResize)
    }
  }, [])

  return (
    <canvas
      ref={canvasRef}
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 0,
        pointerEvents: "none",
      }}
    />
  )
}
// ─── CSS 3D CLOUD ORBS ────────────────────────────────────────────────────────
function Cloud3D() {
  return (
    <div style={{ width:320,height:320,position:"relative",display:"flex",alignItems:"center",justifyContent:"center" }}>
      {/* ambient glow */}
      <div style={{ position:"absolute",inset:0,borderRadius:"50%",
        background:"radial-gradient(circle,rgba(255,107,43,0.18) 0%,transparent 70%)" }}/>
      {/* rings */}
      {[
        { sz:270, dur:8,  rx:60, dir:1,  op:0.7 },
        { sz:300, dur:13, rx:25, dir:-1, op:0.4 },
        { sz:315, dur:20, rx:80, dir:1,  op:0.2 },
      ].map((r,i) => (
        <motion.div key={i}
          style={{ position:"absolute", width:r.sz, height:r.sz, borderRadius:"50%",
            border:`${i===0?2:1}px solid rgba(255,107,43,${r.op})`,
            boxShadow: i===0 ? `0 0 18px rgba(255,107,43,0.25)` : "none",
            transformStyle:"preserve-3d",
          }}
          animate={{ rotateX: r.rx, rotateZ: r.dir===1 ? [0,360] : [360,0] }}
          transition={{ rotateZ:{ duration:r.dur, repeat:Infinity, ease:"linear" } }}
        />
      ))}
      {/* core */}
      <motion.div
        animate={{ scale:[1,1.07,1],
          boxShadow:[
            `0 0 40px rgba(255,107,43,0.55),0 0 80px rgba(255,107,43,0.2),inset 0 0 30px rgba(200,60,0,0.4)`,
            `0 0 65px rgba(255,107,43,0.85),0 0 130px rgba(255,107,43,0.4),inset 0 0 55px rgba(200,60,0,0.6)`,
            `0 0 40px rgba(255,107,43,0.55),0 0 80px rgba(255,107,43,0.2),inset 0 0 30px rgba(200,60,0,0.4)`,
          ]}}
        transition={{ duration:3, repeat:Infinity, ease:"easeInOut" }}
        style={{ width:120,height:120,borderRadius:"50%",position:"relative",zIndex:2,overflow:"hidden",
          background:"radial-gradient(circle at 33% 33%,#FFA060,#FF6B2B 45%,#cc3300 72%,#5a1400)" }}
      >
        {/* scan line */}
        <motion.div animate={{ top:["-10%","115%"] }} transition={{ duration:2.4,repeat:Infinity,ease:"linear" }}
          style={{ position:"absolute",left:0,right:0,height:2,borderRadius:1,
            background:"linear-gradient(90deg,transparent,rgba(255,255,255,0.65),transparent)",pointerEvents:"none" }}/>
      </motion.div>
      {/* orbiting dots */}
      {[0,1,2,3,4].map(i => (
        <motion.div key={i}
          animate={{ rotate:[i*72,i*72+360] }}
          transition={{ duration:5+i*1.5,repeat:Infinity,ease:"linear" }}
          style={{ position:"absolute",width:265,height:265,transformOrigin:"50% 50%" }}>
          <div style={{ position:"absolute",top:0,left:"50%",transform:"translate(-50%,-50%)",
            width:8+i*1.5,height:8+i*1.5,borderRadius:"50%",background:ORANGE,
            boxShadow:`0 0 10px ${ORANGE}`,opacity:0.65+i*0.07 }}/>
        </motion.div>
      ))}
    </div>
  );
}

// ─── GLASS CARD ───────────────────────────────────────────────────────────────
type GlassCardProps = {
  children: React.ReactNode
  style?: React.CSSProperties
  glow?: boolean
}

function GlassCard({ children, style = {}, glow = false }: GlassCardProps) {
  return (
    <div
      style={{
        background: CARD_BG,
        border: `1px solid ${
          glow ? "rgba(255,107,43,0.42)" : CARD_BORDER
        }`,
        borderRadius: 16,
        backdropFilter: "blur(18px)",
        boxShadow: glow
          ? "0 0 32px rgba(255,107,43,0.12), inset 0 1px 0 rgba(255,255,255,0.06)"
          : "inset 0 1px 0 rgba(255,255,255,0.06)",
        ...style,
      }}
    >
      {children}
    </div>
  )
}
// ─── ANIMATED COUNTER ─────────────────────────────────────────────────────────


type AnimatedCounterProps = {
  target: number
  suffix?: string
  duration?: number
}

function AnimatedCounter({
  target,
  suffix = "",
  duration = 2000,
}: AnimatedCounterProps) {
  const [val, setVal] = useState<number>(0)
  const ref = useRef<HTMLSpanElement | null>(null)
  const inView = useInView(ref, { once: true })

  useEffect(() => {
    if (!inView) return

    let raf: number
    const t0 = Date.now()

    const tick = () => {
      const progress = Math.min((Date.now() - t0) / duration, 1)
      const eased = 1 - Math.pow(1 - progress, 3)

      const calculated =
        target % 1 !== 0
          ? parseFloat((target * eased).toFixed(1))
          : Math.round(target * eased)

      setVal(calculated)

      if (progress < 1) {
        raf = requestAnimationFrame(tick)
      }
    }

    raf = requestAnimationFrame(tick)

    return () => cancelAnimationFrame(raf)
  }, [inView, target, duration])

  return (
    <span ref={ref}>
      {val}
      {suffix}
    </span>
  )
}
// ─── CPU GAUGE ───────────────
function CpuGauge({ value = 0 }: { value?: number }) {
  const safeValue = Math.max(0, Math.min(100, Number(value) || 0))

  const dynamicColor =
    safeValue > 80 ? "#FF4A6B" :
    safeValue > 60 ? "#FFB347" :
    ORANGE

  return (
    <div
      style={{
        position: "relative",
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
      }}
    >
      <RadialBarChart
        width={220}
        height={120}
        cx={110}
        cy={110}
        innerRadius={70}
        outerRadius={100}
        startAngle={180}
        endAngle={0}
        data={[{ cpu: safeValue }]}
      >
        <RadialBar
          dataKey="cpu"
          cornerRadius={8}
          fill={dynamicColor}
          background={{ fill: "rgba(255,107,43,0.1)" }}
          isAnimationActive
          animationDuration={800}
        />
      </RadialBarChart>

      <div
        style={{
          marginTop: -58,
          textAlign: "center",
          pointerEvents: "none",
        }}
      >
        <div
          style={{
            fontFamily: "'Courier New', monospace",
            color: dynamicColor,
            fontSize: 36,
            fontWeight: 700,
            textShadow: `0 0 20px ${dynamicColor}`,
          }}
        >
          {safeValue.toFixed(1)}%
        </div>

        <div
          style={{
            color: "rgba(255,255,255,0.42)",
            fontSize: 11,
            letterSpacing: 3,
          }}
        >
          CPU LOAD
        </div>
      </div>
    </div>
  )
}

// ─── STATE BADGE ──────────────────────────────────────────────────────────────
type StateType = keyof typeof stateConfig

type StateBadgeProps = {
  state: StateType | string
}

function StateBadge({ state }: StateBadgeProps) {
  const cfg =
    stateConfig[state as StateType] ?? stateConfig.PENDING

  return (
    <motion.span
      animate={
        state === "EXECUTING"
          ? { opacity: [1, 0.35, 1] }
          : {}
      }
      transition={{ duration: 1.2, repeat: Infinity }}
      style={{
        display: "inline-block",
        padding: "3px 10px",
        borderRadius: 20,
        fontSize: 11,
        fontWeight: 700,
        letterSpacing: 1.5,
        color: cfg.color,
        border: `1px solid ${cfg.color}`,
        boxShadow: `0 0 8px ${cfg.glow}`,
        textTransform: "uppercase",
        fontFamily: "'Courier New',monospace",
      }}
    >
      {state}
    </motion.span>
  )
}

// ─── GRID BG ──────────────────────────────────────────────────────────────────
function GridBg() {
  return (
    <div style={{ position:"fixed",inset:0,zIndex:0,overflow:"hidden",pointerEvents:"none" }}>
      <div style={{ position:"absolute",inset:0,
        backgroundImage:`linear-gradient(rgba(255,107,43,0.05) 1px,transparent 1px),
          linear-gradient(90deg,rgba(255,107,43,0.05) 1px,transparent 1px)`,
        backgroundSize:"60px 60px" }}/>
      <div style={{ position:"absolute",inset:0,
        background:`radial-gradient(ellipse 80% 55% at 50% 0%,rgba(255,107,43,0.09) 0%,transparent 70%)` }}/>
    </div>
  );
}

// ─── NAV ──────────────────────────────────────────────────────────────────────
function Nav() {
  const [active,setActive] = useState(0);
  interface ScrollToIndex {
    (index: number): void;
  }

  const scrollTo: ScrollToIndex = (i: number) => {
    document.getElementById(`sec-${i}`)?.scrollIntoView({ behavior: "smooth" });
    setActive(i);
  };
  return (
    <nav style={{ position:"fixed",top:0,left:0,right:0,zIndex:100,
      background:"rgba(7,10,15,0.88)",backdropFilter:"blur(20px)",
      borderBottom:`1px solid rgba(255,107,43,0.12)`,
      height:60,display:"flex",alignItems:"center",padding:"0 32px",justifyContent:"space-between" }}>
      <div style={{ display:"flex",alignItems:"center",gap:12 }}>
        <motion.div animate={{ rotate:360 }} transition={{ duration:8,repeat:Infinity,ease:"linear" }}
          style={{ width:28,height:28,borderRadius:"50%",border:`2px solid ${ORANGE}`,
            display:"flex",alignItems:"center",justifyContent:"center",boxShadow:`0 0 12px ${ORANGE_GLOW}` }}>
          <div style={{ width:8,height:8,borderRadius:"50%",background:ORANGE }}/>
        </motion.div>
        <span style={{ fontFamily:"'Courier New',monospace",color:WHITE,fontWeight:700,fontSize:16,letterSpacing:3 }}>
          ATLAS <span style={{ color:ORANGE }}>OPS</span>
        </span>
      </div>
      <div style={{ display:"flex",gap:5 }}>
        {navItems.map((item,i)=>(
          <button key={i} onClick={()=>scrollTo(i)} style={{
            background:active===i?`rgba(255,107,43,0.15)`:"transparent",
            border:`1px solid ${active===i?ORANGE:"transparent"}`,
            color:active===i?ORANGE:"rgba(255,255,255,0.42)",
            padding:"5px 13px",borderRadius:6,cursor:"pointer",
            fontSize:11,letterSpacing:1.5,fontFamily:"'Courier New',monospace",transition:"all 0.2s",
          }}>{item.toUpperCase()}</button>
        ))}
      </div>
      <div style={{ display:"flex",alignItems:"center",gap:8 }}>
        <motion.div animate={{ scale:[1,1.5,1] }} transition={{ duration:2,repeat:Infinity }}
          style={{ width:8,height:8,borderRadius:"50%",background:"#4AFFC4",boxShadow:"0 0 8px #4AFFC4" }}/>
        <span style={{ color:"rgba(255,255,255,0.38)",fontSize:11,letterSpacing:2 }}>LIVE</span>
      </div>
    </nav>
  );
}

// ─── CHART TOOLTIP ────────────────────────────────────────────────────────────
type ChartTooltipProps = {
  active?: boolean
  payload?: { value?: number }[]
}

function ChartTooltip({ active, payload }: ChartTooltipProps) {
  if (!active || !payload || payload.length === 0) return null

  return (
    <GlassCard style={{ padding: "10px 16px" }}>
      <div
        style={{
          color: ORANGE,
          fontSize: 12,
          fontFamily: "'Courier New',monospace",
        }}
      >
        CPU: {payload[0]?.value ?? 0}%
      </div>

      <div
        style={{
          color: "#9B8BFF",
          fontSize: 12,
          fontFamily: "'Courier New',monospace",
        }}
      >
        EMA: {payload[1]?.value ?? 0}%
      </div>
    </GlassCard>
  )
}

// ══════════════════════════════════════════════════════════════════════════════
// SECTION 1 — HERO
// ══════════════════════════════════════════════════════════════════════════════


type MonitorData = {
  cpu?: number
  average_cpu?: number
  confidence?: number
  slope?: number
  status?: string
}

function HeroSection({ monitor }: { monitor?: MonitorData }) {
  return (
    <section id="sec-0" style={{ minHeight:"100vh",display:"flex",alignItems:"center",
      justifyContent:"center",padding:"80px 40px 40px",position:"relative",overflow:"hidden" }}>
      <div style={{ position:"absolute",top:"5%",left:"50%",transform:"translateX(-50%)",
        width:700,height:700,borderRadius:"50%",pointerEvents:"none",
        background:"radial-gradient(circle,rgba(255,107,43,0.09) 0%,transparent 70%)" }}/>
      <div style={{ display:"grid",gridTemplateColumns:"1fr 1fr",gap:60,
        alignItems:"center",maxWidth:1200,width:"100%",zIndex:1 }}>
        {/* LEFT */}
        <motion.div initial={{ opacity:0,x:-60 }} animate={{ opacity:1,x:0 }} transition={{ duration:1,ease:[0.16,1,0.3,1] }}>
          <div style={{ fontSize:11,letterSpacing:4,color:ORANGE,fontFamily:"'Courier New',monospace",
            marginBottom:20,display:"flex",alignItems:"center",gap:10 }}>
            <motion.span animate={{ opacity:[1,0,1] }} transition={{ duration:1.5,repeat:Infinity }}>▮</motion.span>
            SYSTEM STATUS · NOMINAL
          </div>
          <h1 style={{ fontSize:"clamp(40px,5.5vw,68px)",fontFamily:"'Courier New',monospace",
            fontWeight:800,color:WHITE,lineHeight:1.05,marginBottom:16,
            textShadow:"0 0 40px rgba(255,255,255,0.07)" }}>
            CLOUD<br/>
            <span style={{ color:ORANGE,textShadow:`0 0 30px ${ORANGE_GLOW},0 0 60px rgba(255,107,43,0.25)` }}>CONTROL</span><br/>
            CENTER
          </h1>
          <p style={{ color:"rgba(255,255,255,0.38)",fontSize:14,lineHeight:1.8,
            maxWidth:420,marginBottom:40,letterSpacing:0.4 }}>
            Autonomous infrastructure intelligence. Real-time anomaly detection,
            adaptive remediation, and predictive scaling — unified.
          </p>
          <div style={{ display: "flex", gap: 16, marginBottom: 40 }}>
  {[
    {
      label: "CPU",
      value: `${(monitor?.average_cpu ?? monitor?.cpu ?? 0).toFixed(1)}%`,
      trend: (monitor?.slope ?? 0) > 0 ? "↑ RISING" : "↓ FALLING",
    },
    {
      label: "CONFIDENCE",
      value: `${Math.round((monitor?.confidence ?? 0) * 100)}%`,
      trend: "↗ STABLE",
    },
    {
      label: "UPTIME",
      value: "99.97%",
      trend: "● LIVE",
    },
  ].map((s, i) => (
    <GlassCard key={i} style={{ padding: "16px 18px", flex: 1 }}>
      <div
        style={{
          fontSize: 10,
          color: "rgba(255,255,255,0.33)",
          letterSpacing: 2,
          marginBottom: 4,
          fontFamily: "'Courier New',monospace",
        }}
      >
        {s.label}
      </div>

      <div
        style={{
          fontSize: 22,
          color: WHITE,
          fontWeight: 700,
          fontFamily: "'Courier New',monospace",
        }}
      >
        {s.value}
      </div>

      <div
        style={{
          fontSize: 10,
          color: ORANGE,
          marginTop: 4,
          letterSpacing: 1,
        }}
      >
        {s.trend}
      </div>
    </GlassCard>
  ))}
</div>
          <motion.button whileHover={{ scale:1.04,boxShadow:`0 0 55px ${ORANGE_GLOW}` }} whileTap={{ scale:0.97 }}
            style={{ background:`linear-gradient(135deg,${ORANGE},${ORANGE_DIM})`,border:"none",
              borderRadius:10,color:WHITE,fontFamily:"'Courier New',monospace",fontWeight:700,
              fontSize:13,letterSpacing:2,padding:"14px 36px",cursor:"pointer",
              position:"relative",overflow:"hidden",
              boxShadow:`0 0 30px ${ORANGE_GLOW},0 4px 20px rgba(0,0,0,0.4)` }}>
            <motion.div animate={{ x:[-200,400] }} transition={{ duration:2.5,repeat:Infinity,ease:"linear" }}
              style={{ position:"absolute",top:0,left:0,width:60,height:"100%",
                background:"linear-gradient(90deg,transparent,rgba(255,255,255,0.25),transparent)",
                transform:"skewX(-20deg)",pointerEvents:"none" }}/>
            ⚡ CREATE INCIDENT
          </motion.button>
        </motion.div>
        {/* RIGHT */}
        <motion.div initial={{ opacity:0,scale:0.82 }} animate={{ opacity:1,scale:1 }}
          transition={{ duration:1.2,ease:[0.16,1,0.3,1],delay:0.2 }}
          style={{ display:"flex",flexDirection:"column",alignItems:"center",gap:24 }}>
          <Cloud3D/>
          <CpuGauge value={73}/>
        </motion.div>
      </div>
    </section>
  );
}

// ══════════════════════════════════════════════════════════════════════════════
// SECTION 2 — ANALYTICS
// ══════════════════════════════════════════════════════════════════════════════

function AnalyticsSection({ monitor }: { monitor: MonitorData }) {

  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref,{once:true,margin:"-100px"});
  return (
    <section id="sec-1" ref={ref} style={{ minHeight:"100vh",display:"flex",
      alignItems:"center",padding:"80px 40px",position:"relative" }}>
      <div style={{ maxWidth:1200,width:"100%",margin:"0 auto" }}>
        <motion.div initial={{ opacity:0,y:30 }} animate={inView?{opacity:1,y:0}:{}} transition={{ duration:0.8 }} style={{ marginBottom:40 }}>
          <div style={{ fontSize:11,color:ORANGE,letterSpacing:4,fontFamily:"'Courier New',monospace",marginBottom:8 }}>02 · TREND ANALYTICS</div>
          <h2 style={{ fontSize:42,fontFamily:"'Courier New',monospace",color:WHITE,fontWeight:700 }}>CPU <span style={{ color:ORANGE }}>Trend</span> Analysis</h2>
        </motion.div>
        <div style={{ display:"grid",gridTemplateColumns:"2fr 1fr",gap:24 }}>
          {/* Chart */}
          <motion.div initial={{ opacity:0,x:-40 }} animate={inView?{opacity:1,x:0}:{}} transition={{ duration:0.9,delay:0.2 }}>
            <GlassCard style={{ padding:"28px 24px" }} glow>
              <div style={{ color:"rgba(255,255,255,0.45)",fontSize:11,letterSpacing:2,marginBottom:20,fontFamily:"'Courier New',monospace" }}>60-SECOND CPU HISTORY · EMA OVERLAY</div>
              <ResponsiveContainer width="100%" height={280}>
                <AreaChart data={cpuHistory}>
                  <defs>
                    <linearGradient id="cpuG" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor={ORANGE} stopOpacity={0.35}/><stop offset="95%" stopColor={ORANGE} stopOpacity={0}/>
                    </linearGradient>
                    <linearGradient id="emaG" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#9B8BFF" stopOpacity={0.2}/><stop offset="95%" stopColor="#9B8BFF" stopOpacity={0}/>
                    </linearGradient>
                  </defs>
                  <CartesianGrid stroke="rgba(255,107,43,0.08)"/>
                  <XAxis dataKey="time" stroke="rgba(255,255,255,0.12)" tick={{fill:"rgba(255,255,255,0.28)",fontSize:10}} interval={9}/>
                  <YAxis stroke="rgba(255,255,255,0.12)" tick={{fill:"rgba(255,255,255,0.28)",fontSize:10}} domain={[0,100]}/>
                  <Tooltip content={<ChartTooltip/>}/>
                  <Area type="monotone" dataKey="cpu" stroke={ORANGE} strokeWidth={2} fill="url(#cpuG)" dot={false}/>
                  <Area type="monotone" dataKey="ema" stroke="#9B8BFF" strokeWidth={1.5} fill="url(#emaG)" strokeDasharray="5 3" dot={false}/>
                </AreaChart>
              </ResponsiveContainer>
              <div style={{ display:"flex",gap:24,marginTop:16 }}>
                {[{color:ORANGE,label:"CPU Usage"},{color:"#9B8BFF",label:"EMA Trend"}].map((l,i)=>(
                  <div key={i} style={{ display:"flex",alignItems:"center",gap:8 }}>
                    <div style={{ width:24,height:2,background:l.color,borderRadius:2 }}/>
                    <span style={{ color:"rgba(255,255,255,0.38)",fontSize:11,letterSpacing:1 }}>{l.label}</span>
                  </div>
                ))}
              </div>
            </GlassCard>
          </motion.div>
          {/* Side */}
          <div style={{ display:"flex",flexDirection:"column",gap:20 }}>
            {[0.3,0.4,0.5].map((_,i)=>(
              <motion.div key={i} initial={{ opacity:0,x:40 }} animate={inView?{opacity:1,x:0}:{}} transition={{ duration:0.9,delay:_}}>
                {i===0&&(
                  <GlassCard style={{ padding:"24px" }}>
                    <div style={{ fontSize:11,color:"rgba(255,255,255,0.32)",letterSpacing:2,marginBottom:16,fontFamily:"'Courier New',monospace" }}>CONFIDENCE SCORE</div>
                    <ResponsiveContainer width="100%" height={110}>
                      <RadialBarChart cx="50%" cy="100%" innerRadius="60%" outerRadius="90%" startAngle={180} endAngle={0} data={[{value:87,fill:ORANGE}]}>
                        <RadialBar background={{fill:"rgba(255,107,43,0.1)"}} dataKey="value" cornerRadius={8}/>
                      </RadialBarChart>
                    </ResponsiveContainer>
                    <div style={{ textAlign:"center",marginTop:-36 }}>
                      <div style={{ fontSize:34,color:ORANGE,fontWeight:800,fontFamily:"'Courier New',monospace",textShadow:`0 0 20px ${ORANGE_GLOW}` }}>87%</div>
                      <div style={{ fontSize:10,color:"rgba(255,255,255,0.32)",letterSpacing:2 }}>HIGH CONFIDENCE</div>
                    </div>
                  </GlassCard>
                )}
                {i===1&&(
                  <GlassCard style={{ padding:"24px" }}>
                    <div style={{ fontSize:11,color:"rgba(255,255,255,0.32)",letterSpacing:2,marginBottom:16,fontFamily:"'Courier New',monospace" }}>SLOPE DIRECTION</div>
                    <div style={{ display:"flex",alignItems:"center",gap:14 }}>
                      <motion.div animate={{ y:[-4,4,-4] }} transition={{ duration:2,repeat:Infinity }}
                        style={{ fontSize:48,color:ORANGE,textShadow:`0 0 20px ${ORANGE_GLOW}`,lineHeight:1 }}>↗</motion.div>
                      <div>
                        <div style={{ color:WHITE,fontSize:22,fontWeight:700,fontFamily:"'Courier New',monospace" }}>+2.4%</div>
                        <div style={{ color:"rgba(255,255,255,0.32)",fontSize:11,letterSpacing:1 }}>PER MINUTE</div>
                        <div style={{ color:ORANGE,fontSize:11,marginTop:4,letterSpacing:1 }}>ESCALATING</div>
                      </div>
                    </div>
                  </GlassCard>
                )}
                {i===2&&(
                  <GlassCard style={{ padding:"20px",borderColor:"rgba(255,74,107,0.35)" }}>
                    <motion.div animate={{ opacity:[1,0.35,1] }} transition={{ duration:1.5,repeat:Infinity }}
                      style={{ display:"flex",alignItems:"center",gap:10 }}>
                      <div style={{ width:10,height:10,borderRadius:"50%",background:"#FF4A6B",
                        boxShadow:"0 0 10px rgba(255,74,107,0.9)",flexShrink:0 }}/>
                      <div style={{ color:"#FF4A6B",fontSize:11,fontWeight:700,letterSpacing:2,fontFamily:"'Courier New',monospace" }}>THRESHOLD BREACH</div>
                    </motion.div>
                    <div style={{ color:"rgba(255,255,255,0.42)",fontSize:12,marginTop:10,lineHeight:1.6 }}>
                      CPU exceeds 70% threshold. Auto-remediation triggered.
                    </div>
                  </GlassCard>
                )}
              </motion.div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

// ══════════════════════════════════════════════════════════════════════════════
// SECTION 3 — INCIDENTS
// ══════════════════════════════════════════════════════════════════════════════
function IncidentsSection() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref,{once:true,margin:"-100px"});
  return (
    <section id="sec-2" ref={ref} style={{ minHeight:"100vh",display:"flex",alignItems:"center",padding:"80px 40px",position:"relative" }}>
      <div style={{ maxWidth:1200,width:"100%",margin:"0 auto" }}>
        <motion.div initial={{ opacity:0,y:30 }} animate={inView?{opacity:1,y:0}:{}} transition={{ duration:0.8 }} style={{ marginBottom:40 }}>
          <div style={{ fontSize:11,color:ORANGE,letterSpacing:4,fontFamily:"'Courier New',monospace",marginBottom:8 }}>03 · CONTROL PANEL</div>
          <h2 style={{ fontSize:42,fontFamily:"'Courier New',monospace",color:WHITE,fontWeight:700 }}>Incident <span style={{ color:ORANGE }}>Queue</span></h2>
        </motion.div>
        <motion.div initial={{ opacity:0,y:40 }} animate={inView?{opacity:1,y:0}:{}} transition={{ duration:0.9,delay:0.2 }}>
          <GlassCard glow>
            <div style={{ display:"grid",gridTemplateColumns:"1.5fr 1.2fr 0.8fr 1fr 1fr",
              padding:"16px 24px",borderBottom:`1px solid rgba(255,107,43,0.12)` }}>
              {["INCIDENT ID","STATE","ATTEMPTS","CONFIDENCE","ACTION"].map((h,i)=>(
                <div key={i} style={{ fontSize:10,color:"rgba(255,255,255,0.28)",letterSpacing:2,fontFamily:"'Courier New',monospace" }}>{h}</div>
              ))}
            </div>
            {incidents.map((inc,i)=>(
              <motion.div key={inc.id}
                initial={{ opacity:0,x:-20 }} animate={inView?{opacity:1,x:0}:{}}
                transition={{ duration:0.5,delay:0.1*i+0.3 }}
                whileHover={{ background:"rgba(255,107,43,0.05)" }}
                style={{ display:"grid",gridTemplateColumns:"1.5fr 1.2fr 0.8fr 1fr 1fr",
                  padding:"18px 24px",borderBottom:`1px solid rgba(255,255,255,0.04)`,
                  alignItems:"center",transition:"background 0.2s" }}>
                <div style={{ fontFamily:"'Courier New',monospace",color:WHITE,fontSize:13,fontWeight:600 }}>{inc.id}</div>
                <div><StateBadge state={inc.state}/></div>
                <div style={{ fontFamily:"'Courier New',monospace",color:"rgba(255,255,255,0.62)",fontSize:14 }}>{inc.attempts}x</div>
                <div style={{ display:"flex",alignItems:"center",gap:10 }}>
                  <div style={{ flex:1,height:4,background:"rgba(255,255,255,0.08)",borderRadius:2,overflow:"hidden" }}>
                    <motion.div initial={{ width:0 }} animate={inView?{width:`${inc.confidence}%`}:{width:0}}
                      transition={{ duration:1,delay:0.2*i+0.5 }}
                      style={{ height:"100%",background:ORANGE,borderRadius:2 }}/>
                  </div>
                  <span style={{ fontSize:12,color:ORANGE,fontFamily:"'Courier New',monospace",minWidth:36 }}>{inc.confidence}%</span>
                </div>
                <motion.button whileHover={{ scale:1.07,boxShadow:`0 0 20px ${ORANGE_GLOW}` }} whileTap={{ scale:0.95 }}
                  disabled={inc.state==="RESOLVED"}
                  style={{ background:inc.state==="RESOLVED"?"rgba(255,255,255,0.04)":`rgba(255,107,43,0.15)`,
                    border:`1px solid ${inc.state==="RESOLVED"?"rgba(255,255,255,0.08)":ORANGE}`,
                    color:inc.state==="RESOLVED"?"rgba(255,255,255,0.22)":ORANGE,
                    borderRadius:6,padding:"6px 16px",cursor:inc.state==="RESOLVED"?"default":"pointer",
                    fontSize:11,letterSpacing:1.5,fontFamily:"'Courier New',monospace",
                    fontWeight:700,transition:"all 0.2s",width:"fit-content" }}>
                  {inc.state==="RESOLVED"?"DONE":"APPROVE"}
                </motion.button>
              </motion.div>
            ))}
          </GlassCard>
        </motion.div>
      </div>
    </section>
  );
}

// ══════════════════════════════════════════════════════════════════════════════
// SECTION 4 — TIMELINE
// ══════════════════════════════════════════════════════════════════════════════
// ─────────────────────────────────────────────────────────────
// TIMELINE SECTION (SAFE + UNIQUE NAMES)
// ─────────────────────────────────────────────────────────────

type ExecStatus = keyof typeof statusColors

type ExecTimelineItem = {
  action: string
  time: string
  cpuBefore: number
  cpuAfter: number
  status: ExecStatus
}

const executionTimeline: ExecTimelineItem[] = [
  {
    action: "Scale Up Instance",
    time: "14:02:31",
    cpuBefore: 86,
    cpuAfter: 42,
    status: "success",
  },
  {
    action: "Warm Pool Activated",
    time: "14:02:45",
    cpuBefore: 78,
    cpuAfter: 55,
    status: "partial",
  },
  {
    action: "Load Rebalanced",
    time: "14:03:02",
    cpuBefore: 71,
    cpuAfter: 49,
    status: "success",
  },
]

function TimelineSection() {
  const ref = useRef<HTMLDivElement | null>(null)
  const inView = useInView(ref, { once: true, margin: "-100px" })

  return (
    <section
      id="sec-3"
      ref={ref}
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        padding: "80px 40px",
        position: "relative",
      }}
    >
      <div style={{ maxWidth: 900, width: "100%", margin: "0 auto" }}>
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.8 }}
          style={{ marginBottom: 60 }}
        >
          <div
            style={{
              fontSize: 11,
              color: ORANGE,
              letterSpacing: 4,
              fontFamily: "'Courier New',monospace",
              marginBottom: 8,
            }}
          >
            04 · EXECUTION LOG
          </div>

          <h2
            style={{
              fontSize: 42,
              fontFamily: "'Courier New',monospace",
              color: WHITE,
              fontWeight: 700,
            }}
          >
            Remediation{" "}
            <span style={{ color: ORANGE }}>Timeline</span>
          </h2>
        </motion.div>

        <div style={{ position: "relative" }}>
          <motion.div
            initial={{ scaleY: 0 }}
            animate={inView ? { scaleY: 1 } : {}}
            transition={{ duration: 1.5, ease: "easeInOut" }}
            style={{
              position: "absolute",
              left: 24,
              top: 0,
              bottom: 0,
              width: 2,
              background: `linear-gradient(to bottom,${ORANGE},rgba(255,107,43,0.06))`,
              transformOrigin: "top",
            }}
          />

          {executionTimeline.map((item, i) => {
            const color = statusColors[item.status]

            return (
              <motion.div
                key={i}
                initial={{ opacity: 0, x: -40 }}
                animate={inView ? { opacity: 1, x: 0 } : {}}
                transition={{ duration: 0.7, delay: i * 0.15 + 0.3 }}
                style={{
                  display: "flex",
                  alignItems: "flex-start",
                  gap: 32,
                  marginBottom: 24,
                  position: "relative",
                  paddingLeft: 64,
                }}
              >
                <motion.div
                  initial={{ scale: 0 }}
                  animate={inView ? { scale: 1 } : {}}
                  transition={{
                    duration: 0.4,
                    delay: i * 0.15 + 0.5,
                    type: "spring",
                  }}
                  style={{
                    position: "absolute",
                    left: 14,
                    top: 18,
                    width: 22,
                    height: 22,
                    borderRadius: "50%",
                    background: color,
                    boxShadow: `0 0 16px ${color}99`,
                    border: `2px solid ${DARK_BG}`,
                    zIndex: 1,
                  }}
                />

                <GlassCard style={{ flex: 1, padding: "20px 24px" }}>
                  <div
                    style={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "flex-start",
                      marginBottom: 12,
                    }}
                  >
                    <div
                      style={{
                        color: WHITE,
                        fontWeight: 700,
                        fontSize: 15,
                      }}
                    >
                      {item.action}
                    </div>

                    <div
                      style={{
                        fontFamily: "'Courier New',monospace",
                        color: "rgba(255,255,255,0.28)",
                        fontSize: 12,
                      }}
                    >
                      {item.time}
                    </div>
                  </div>

                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: 18,
                    }}
                  >
                    <div style={{ display: "flex", gap: 8 }}>
                      <span style={{ fontSize: 11, color: "rgba(255,255,255,0.32)" }}>
                        BEFORE
                      </span>
                      <span
                        style={{
                          fontFamily: "'Courier New',monospace",
                          color: "#FF4A6B",
                          fontSize: 16,
                          fontWeight: 700,
                        }}
                      >
                        {item.cpuBefore}%
                      </span>
                    </div>

                    <motion.span
                      animate={{ x: [0, 5, 0] }}
                      transition={{ duration: 1.5, repeat: Infinity }}
                      style={{ color: ORANGE, fontSize: 18 }}
                    >
                      →
                    </motion.span>

                    <div style={{ display: "flex", gap: 8 }}>
                      <span style={{ fontSize: 11, color: "rgba(255,255,255,0.32)" }}>
                        AFTER
                      </span>
                      <span
                        style={{
                          fontFamily: "'Courier New',monospace",
                          color: "#4AFFC4",
                          fontSize: 16,
                          fontWeight: 700,
                        }}
                      >
                        {item.cpuAfter}%
                      </span>
                    </div>

                    <span
                      style={{
                        marginLeft: "auto",
                        fontSize: 10,
                        padding: "3px 10px",
                        borderRadius: 20,
                        color: color,
                        border: `1px solid ${color}`,
                        boxShadow: `0 0 8px ${color}44`,
                        textTransform: "uppercase",
                        letterSpacing: 1.5,
                        fontFamily: "'Courier New',monospace",
                      }}
                    >
                      {item.status}
                    </span>
                  </div>
                </GlassCard>
              </motion.div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
// ══════════════════════════════════════════════════════════════════════════════
// SECTION 5 — DLQ
// ══════════════════════════════════════════════════════════════════════════════
function DLQSection() {
  const ref = useRef<HTMLDivElement>(null);
  const inView = useInView(ref,{once:true,margin:"-100px"});
  return (
    <section id="sec-4" ref={ref} style={{ minHeight:"100vh",display:"flex",alignItems:"center",padding:"80px 40px",position:"relative" }}>
      <div style={{ maxWidth:1000,width:"100%",margin:"0 auto" }}>
        <motion.div initial={{ opacity:0,y:30 }} animate={inView?{opacity:1,y:0}:{}} transition={{ duration:0.8 }} style={{ marginBottom:40 }}>
          <div style={{ fontSize:11,color:ORANGE,letterSpacing:4,fontFamily:"'Courier New',monospace",marginBottom:8 }}>05 · DEAD LETTER QUEUE</div>
          <h2 style={{ fontSize:42,fontFamily:"'Courier New',monospace",color:WHITE,fontWeight:700 }}>DLQ <span style={{ color:"#FF4A6B" }}>Monitor</span></h2>
        </motion.div>
        <motion.div initial={{ opacity:0,scale:0.95 }} animate={inView?{opacity:1,scale:1}:{}} transition={{ duration:0.8,delay:0.2 }}>
          <div style={{ background:"rgba(255,74,107,0.07)",border:"1px solid rgba(255,74,107,0.4)",
            borderRadius:16,padding:"24px 32px",marginBottom:24,
            display:"flex",alignItems:"center",gap:24,backdropFilter:"blur(16px)" }}>
            <motion.div animate={{ scale:[1,1.3,1],opacity:[1,0.5,1] }} transition={{ duration:1.5,repeat:Infinity }}
              style={{ width:60,height:60,borderRadius:"50%",flexShrink:0,
                background:"rgba(255,74,107,0.12)",border:"2px solid #FF4A6B",
                display:"flex",alignItems:"center",justifyContent:"center",
                boxShadow:"0 0 30px rgba(255,74,107,0.4)",fontSize:26 }}>⚠</motion.div>
            <div style={{ flex:1 }}>
              <div style={{ color:"#FF4A6B",fontFamily:"'Courier New',monospace",fontWeight:700,fontSize:18,marginBottom:6 }}>
                {dlqItems.length} FAILED MESSAGES IN QUEUE
              </div>
              <div style={{ color:"rgba(255,255,255,0.42)",fontSize:13,lineHeight:1.6 }}>
                These incidents could not be processed and require manual intervention.
              </div>
            </div>
            <motion.button whileHover={{ scale:1.05 }} whileTap={{ scale:0.95 }}
              style={{ background:"rgba(255,74,107,0.12)",border:"1px solid #FF4A6B",
                color:"#FF4A6B",borderRadius:8,padding:"10px 20px",cursor:"pointer",
                fontFamily:"'Courier New',monospace",fontSize:12,fontWeight:700,letterSpacing:1.5 }}>RETRY ALL</motion.button>
          </div>
          {dlqItems.map((item,i)=>(
            <motion.div key={item.id} initial={{ opacity:0,x:-30 }} animate={inView?{opacity:1,x:0}:{}}
              transition={{ duration:0.6,delay:i*0.15+0.4 }} style={{ marginBottom:16 }}>
              <GlassCard style={{ padding:"20px 24px",borderColor:"rgba(255,74,107,0.2)" }}>
                <div style={{ display:"flex",justifyContent:"space-between",alignItems:"center" }}>
                  <div style={{ display:"flex",alignItems:"center",gap:16 }}>
                    <motion.div animate={{ opacity:[1,0.15,1] }} transition={{ duration:2,repeat:Infinity,delay:i*0.5 }}
                      style={{ width:8,height:8,borderRadius:"50%",background:"#FF4A6B",
                        boxShadow:"0 0 10px rgba(255,74,107,0.9)",flexShrink:0 }}/>
                    <span style={{ fontFamily:"'Courier New',monospace",color:WHITE,fontWeight:700,fontSize:14 }}>{item.id}</span>
                    <span style={{ color:"rgba(255,255,255,0.38)",fontSize:13 }}>{item.error}</span>
                  </div>
                  <div style={{ display:"flex",alignItems:"center",gap:16 }}>
                    <span style={{ fontFamily:"'Courier New',monospace",color:"rgba(255,255,255,0.22)",fontSize:12 }}>{item.ts}</span>
                    <motion.button whileHover={{ scale:1.05,borderColor:ORANGE,color:ORANGE }}
                      style={{ background:"transparent",border:"1px solid rgba(255,107,43,0.35)",
                        color:"rgba(255,107,43,0.65)",borderRadius:6,padding:"5px 14px",
                        cursor:"pointer",fontSize:11,fontFamily:"'Courier New',monospace",
                        letterSpacing:1,transition:"all 0.2s" }}>RETRY</motion.button>
                  </div>
                </div>
              </GlassCard>
            </motion.div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}

// ══════════════════════════════════════════════════════════════════════════════
// SECTION 6 — METRICS
// ══════════════════════════════════════════════════════════════════════════════
// ─────────────────────────────────────────────────────────────
// METRICS SECTION (STRICT SAFE)
// ─────────────────────────────────────────────────────────────
type MetricsCardType = {
  label: string
  value: number
  suffix: string
  color: string
}

const performanceMetricCards: MetricsCardType[] = [
  { label: "CPU AVG", value: 64, suffix: "%", color: "#FF6B2B" },
  { label: "ANOMALIES", value: 12, suffix: "", color: "#FF4A6B" },
  { label: "AUTOMATIONS", value: 87, suffix: "%", color: "#4AFFC4" },
  { label: "INCIDENTS", value: 6, suffix: "", color: "#9B8BFF" },
]

function MetricsSection() {
const ref = useRef<HTMLDivElement | null>(null)

const inView = useInView(ref, { once: true, margin: "-100px" })

const sparkData = Array.from({ length: 20 }).map((_, i) => ({
  v: 55 + Math.sin(i * 0.8) * 25 + (i % 3) * 3,
}))
  return (
    <section
      id="sec-5"
      ref={ref}
      style={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        padding: "80px 40px",
        position: "relative",
      }}
    >
      <div style={{ maxWidth: 1200, width: "100%", margin: "0 auto" }}>
        <motion.div
          initial={{ opacity: 0, y: 30 }}
          animate={inView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.8 }}
          style={{ marginBottom: 60 }}
        >
          <div
            style={{
              fontSize: 11,
              color: ORANGE,
              letterSpacing: 4,
              fontFamily: "'Courier New',monospace",
              marginBottom: 8,
            }}
          >
            06 · METRICS SUMMARY
          </div>

          <h2
            style={{
              fontSize: 42,
              fontFamily: "'Courier New',monospace",
              color: WHITE,
              fontWeight: 700,
            }}
          >
            System <span style={{ color: ORANGE }}>Performance</span>
          </h2>
        </motion.div>

        <div
          style={{
            display: "grid",
            gridTemplateColumns: "repeat(4,1fr)",
            gap: 20,
            marginBottom: 28,
          }}
        >
          {performanceMetricCards.map((m, i) => (
            <motion.div
              key={i}
              initial={{ opacity: 0, y: 30 }}
              animate={inView ? { opacity: 1, y: 0 } : {}}
              transition={{ duration: 0.7, delay: i * 0.1 + 0.2 }}
            >
              <GlassCard
                style={{ padding: "28px 20px", textAlign: "center" }}
                glow={i === 0}
              >
                <motion.div
                  animate={{
                    textShadow: [
                      `0 0 0px ${m.color}00`,
                      `0 0 28px ${m.color}88`,
                      `0 0 0px ${m.color}00`,
                    ],
                  }}
                  transition={{ duration: 3, repeat: Infinity, delay: i * 0.5 }}
                  style={{
                    fontSize: 46,
                    fontFamily: "'Courier New',monospace",
                    fontWeight: 800,
                    color: m.color,
                    lineHeight: 1,
                    marginBottom: 8,
                  }}
                >
                  {inView && (
                    <AnimatedCounter
                      target={m.value}
                      suffix={m.suffix}
                      duration={2000 + i * 300}
                    />
                  )}
                </motion.div>

                <div
                  style={{
                    color: "rgba(255,255,255,0.32)",
                    fontSize: 11,
                    letterSpacing: 2,
                    fontFamily: "'Courier New',monospace",
                    marginBottom: 16,
                  }}
                >
                  {m.label.toUpperCase()}
                </div>

                <div style={{ opacity: 0.55 }}>
                  <ResponsiveContainer width="100%" height={36}>
                    <LineChart data={sparkData}>
                      <Line
                        type="monotone"
                        dataKey="v"
                        stroke={m.color}
                        strokeWidth={1.5}
                        dot={false}
                      />
                    </LineChart>
                  </ResponsiveContainer>
                </div>
              </GlassCard>
            </motion.div>
          ))}
        </div>
      </div>
    </section>
  )
}

// ══════════════════════════════════════════════════════════════════════════════
// ROOT
// ══════════════════════════════════════════════════════════════════════════════


export default function AtlasOps() {
  const [monitor, setMonitor] = useState<MonitorData | null>(null)
  const [error, setError] = useState<boolean>(false)

  useEffect(() => {
    let isMounted = true
    let isFetching = false

    const fetchData = async () => {
      if (isFetching) return
      isFetching = true

      try {
        setError(false)

        const controller = new AbortController()
        const timeout = setTimeout(() => controller.abort(), 5000)

        const res = await fetch("http://localhost:8081/monitor", {
          signal: controller.signal,
        })

        clearTimeout(timeout)

        if (!res.ok) {
          throw new Error("Backend error")
        }

        const json = await res.json()

        if (isMounted) {
          setMonitor(json.data as MonitorData)
        }
      } catch (err) {
        console.error("Monitor fetch failed:", err)
        if (isMounted) setError(true)
      } finally {
        isFetching = false
      }
    }

    fetchData()
    const interval = setInterval(fetchData, 10000)

    return () => {
      isMounted = false
      clearInterval(interval)
    }
  }, [])

  if (error) {
    return (
      <div style={{ padding: 40, color: "red" }}>
        Backend connection failed
      </div>
    )
  }

  if (!monitor) {
    return (
      <div style={{ padding: 40, color: "orange" }}>
        Connecting to backend...
      </div>
    )
  }

  return (
    <div
      style={{
        background: DARK_BG,
        color: WHITE,
        minHeight: "100vh",
        fontFamily: "system-ui,-apple-system,sans-serif",
        overflowX: "hidden",
        position: "relative",
      }}
    >
      <GridBg />
      <ParticleCanvas />
      <Nav />
      <HeroSection monitor={monitor} />
      <AnalyticsSection monitor={monitor} />
      <IncidentsSection />
      <TimelineSection />
      <DLQSection />
      <MetricsSection />
    </div>
  )
}