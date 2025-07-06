import { useState, useRef } from 'react'
import { DrawingCanvas, DrawingCanvasRef } from './components/DrawingCanvas'
import { MagicalEffects } from './components/MagicalEffects'
import { MagicWand } from './components/MagicWand'
import { KeyboardShortcuts } from './components/KeyboardShortcuts'
import { useKeyboardShortcuts } from './hooks/useKeyboardShortcuts'
import './App.css'

function App() {
  const [brushColor, setBrushColor] = useState('#FF6B9D')
  const [brushSize, setBrushSize] = useState(8)
  const [magicMode, setMagicMode] = useState('normal')
  const canvasRef = useRef<DrawingCanvasRef>(null)

  const colors = [
    '#FF6B9D', // Pink
    '#4ECDC4', // Teal
    '#45B7D1', // Blue
    '#96CEB4', // Green
    '#FFEAA7', // Yellow
    '#DDA0DD', // Plum
    '#FFB347', // Orange
    '#FF6B6B', // Red
  ]

  const handleMagicCast = (effect: string) => {
    setMagicMode(effect)
    // Reset to normal after 5 seconds
    setTimeout(() => setMagicMode('normal'), 5000)
  }

  const handleClearCanvas = () => {
    canvasRef.current?.clearCanvas()
  }

  // Set up keyboard shortcuts
  useKeyboardShortcuts({
    colors,
    brushSize,
    setBrushColor,
    setBrushSize,
    onClearCanvas: handleClearCanvas,
    onMagicCast: handleMagicCast,
  })

  return (
    <MagicalEffects>
      <div style={{ 
        padding: '20px', 
        backgroundColor: '#F8F3FF',
        minHeight: '100vh',
        fontFamily: 'inherit'
      }}>
        <h1 style={{ 
          textAlign: 'center', 
          color: '#6A4C93',
          fontSize: '3rem',
          fontWeight: '700',
          textShadow: '2px 2px 4px rgba(0,0,0,0.1)',
          animation: magicMode !== 'normal' ? 'pulse 1s infinite' : 'none'
        }}>
          🎨 PlaySlate 🌟
        </h1>
      
      <div style={{ 
        display: 'flex', 
        justifyContent: 'center', 
        gap: '20px',
        marginBottom: '20px',
        flexWrap: 'wrap'
      }}>
        <div style={{ textAlign: 'center' }}>
          <h3 style={{ color: '#6A4C93', margin: '0 0 10px 0', fontWeight: '600' }}>Colors</h3>
          <div style={{ display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
            {colors.map((color) => (
              <button
                key={color}
                onClick={() => setBrushColor(color)}
                style={{
                  width: '50px',
                  height: '50px',
                  backgroundColor: color,
                  border: brushColor === color ? '4px solid #333' : '2px solid #fff',
                  borderRadius: '50%',
                  cursor: 'pointer',
                  boxShadow: '0 4px 8px rgba(0,0,0,0.2)',
                  transform: brushColor === color ? 'scale(1.1)' : 'scale(1)',
                  transition: 'all 0.2s ease',
                }}
              />
            ))}
          </div>
        </div>

        <div style={{ textAlign: 'center' }}>
          <h3 style={{ color: '#6A4C93', margin: '0 0 10px 0', fontWeight: '600' }}>Brush Size</h3>
          <input
            type="range"
            min="2"
            max="50"
            value={brushSize}
            onChange={(e) => setBrushSize(Number(e.target.value))}
            style={{
              width: '200px',
              height: '20px',
              borderRadius: '10px',
              outline: 'none',
            }}
          />
          <div style={{ 
            marginTop: '10px',
            width: `${brushSize}px`,
            height: `${brushSize}px`,
            backgroundColor: brushColor,
            borderRadius: '50%',
            margin: '10px auto',
            border: '2px solid #fff',
            boxShadow: '0 2px 4px rgba(0,0,0,0.2)'
          }} />
        </div>
      </div>

      <MagicWand onMagicCast={handleMagicCast} />

      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <DrawingCanvas 
          ref={canvasRef}
          width={800} 
          height={600} 
          brushColor={magicMode === 'rainbow' ? '#FF0000' : brushColor}
          brushSize={brushSize}
        />
      </div>
      
      <KeyboardShortcuts />
      </div>
    </MagicalEffects>
  )
}

export default App
