import React, { useRef, useEffect, useState, useCallback, forwardRef, useImperativeHandle } from 'react';

interface Point {
  x: number;
  y: number;
}

interface DrawingCanvasProps {
  width?: number;
  height?: number;
  brushColor?: string;
  brushSize?: number;
}

export interface DrawingCanvasRef {
  clearCanvas: () => void;
}

export const DrawingCanvas = forwardRef<DrawingCanvasRef, DrawingCanvasProps>(({
  width = 800,
  height = 600,
  brushColor = '#FF6B9D',
  brushSize = 8,
}, ref) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);
  const [lastPoint, setLastPoint] = useState<Point | null>(null);

  const getCanvasCoordinates = useCallback((event: MouseEvent | TouchEvent): Point => {
    const canvas = canvasRef.current;
    if (!canvas) return { x: 0, y: 0 };

    const rect = canvas.getBoundingClientRect();
    const scaleX = canvas.width / rect.width;
    const scaleY = canvas.height / rect.height;

    if ('touches' in event) {
      const touch = event.touches[0] || event.changedTouches[0];
      return {
        x: (touch.clientX - rect.left) * scaleX,
        y: (touch.clientY - rect.top) * scaleY,
      };
    } else {
      return {
        x: (event.clientX - rect.left) * scaleX,
        y: (event.clientY - rect.top) * scaleY,
      };
    }
  }, []);

  const drawLine = useCallback((from: Point, to: Point) => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext('2d');
    if (!ctx) return;

    ctx.lineCap = 'round';
    ctx.lineJoin = 'round';
    ctx.strokeStyle = brushColor;
    ctx.lineWidth = brushSize;

    ctx.beginPath();
    ctx.moveTo(from.x, from.y);
    ctx.lineTo(to.x, to.y);
    ctx.stroke();
  }, [brushColor, brushSize]);

  const startDrawing = useCallback((event: MouseEvent | TouchEvent) => {
    event.preventDefault();
    const point = getCanvasCoordinates(event);
    setIsDrawing(true);
    setLastPoint(point);

    // Draw a dot for single clicks/taps
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext('2d');
    if (ctx) {
      ctx.beginPath();
      ctx.arc(point.x, point.y, brushSize / 2, 0, 2 * Math.PI);
      ctx.fillStyle = brushColor;
      ctx.fill();
    }
  }, [getCanvasCoordinates, brushColor, brushSize]);

  const draw = useCallback((event: MouseEvent | TouchEvent) => {
    if (!isDrawing || !lastPoint) return;
    
    event.preventDefault();
    const currentPoint = getCanvasCoordinates(event);
    drawLine(lastPoint, currentPoint);
    setLastPoint(currentPoint);
  }, [isDrawing, lastPoint, getCanvasCoordinates, drawLine]);

  const stopDrawing = useCallback(() => {
    setIsDrawing(false);
    setLastPoint(null);
  }, []);

  useEffect(() => {
    const canvas = canvasRef.current;
    if (!canvas) return;

    // Mouse events
    const handleMouseDown = (e: MouseEvent) => startDrawing(e);
    const handleMouseMove = (e: MouseEvent) => draw(e);
    const handleMouseUp = () => stopDrawing();

    // Touch events
    const handleTouchStart = (e: TouchEvent) => startDrawing(e);
    const handleTouchMove = (e: TouchEvent) => draw(e);
    const handleTouchEnd = () => stopDrawing();

    canvas.addEventListener('mousedown', handleMouseDown);
    canvas.addEventListener('mousemove', handleMouseMove);
    canvas.addEventListener('mouseup', handleMouseUp);
    canvas.addEventListener('mouseleave', handleMouseUp);

    canvas.addEventListener('touchstart', handleTouchStart);
    canvas.addEventListener('touchmove', handleTouchMove);
    canvas.addEventListener('touchend', handleTouchEnd);

    return () => {
      canvas.removeEventListener('mousedown', handleMouseDown);
      canvas.removeEventListener('mousemove', handleMouseMove);
      canvas.removeEventListener('mouseup', handleMouseUp);
      canvas.removeEventListener('mouseleave', handleMouseUp);

      canvas.removeEventListener('touchstart', handleTouchStart);
      canvas.removeEventListener('touchmove', handleTouchMove);
      canvas.removeEventListener('touchend', handleTouchEnd);
    };
  }, [startDrawing, draw, stopDrawing]);

  const clearCanvas = useCallback(() => {
    const canvas = canvasRef.current;
    const ctx = canvas?.getContext('2d');
    if (ctx && canvas) {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      
      // Set a magical background gradient
      const gradient = ctx.createLinearGradient(0, 0, canvas.width, canvas.height);
      gradient.addColorStop(0, '#FFE5F1');
      gradient.addColorStop(1, '#E5F3FF');
      ctx.fillStyle = gradient;
      ctx.fillRect(0, 0, canvas.width, canvas.height);
    }
  }, []);

  useEffect(() => {
    clearCanvas();
  }, [clearCanvas]);

  // Expose clearCanvas method via ref
  useImperativeHandle(ref, () => ({
    clearCanvas,
  }));

  return (
    <div className="drawing-canvas-container">
      <canvas
        ref={canvasRef}
        width={width}
        height={height}
        style={{
          border: '4px solid #FFB3E6',
          borderRadius: '20px',
          cursor: 'crosshair',
          touchAction: 'none',
        }}
      />
      <button
        onClick={clearCanvas}
        style={{
          marginTop: '10px',
          padding: '12px 24px',
          fontSize: '18px',
          backgroundColor: '#FF6B9D',
          color: 'white',
          border: 'none',
          borderRadius: '25px',
          cursor: 'pointer',
          fontWeight: 'bold',
        }}
      >
        ✨ Clear Canvas
      </button>
    </div>
  );
});

DrawingCanvas.displayName = 'DrawingCanvas';