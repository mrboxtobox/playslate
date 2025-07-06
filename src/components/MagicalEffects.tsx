import React, { useEffect, useState } from 'react';

interface Sparkle {
  id: number;
  x: number;
  y: number;
  size: number;
  color: string;
  life: number;
  maxLife: number;
}

interface MagicalEffectsProps {
  children: React.ReactNode;
}

export const MagicalEffects: React.FC<MagicalEffectsProps> = ({ children }) => {
  const [sparkles, setSparkles] = useState<Sparkle[]>([]);

  const colors = ['#FFD700', '#FF69B4', '#00BFFF', '#98FB98', '#DDA0DD', '#FFB347'];

  const createSparkle = (x: number, y: number): Sparkle => ({
    id: Math.random(),
    x: x + (Math.random() - 0.5) * 100,
    y: y + (Math.random() - 0.5) * 100,
    size: Math.random() * 6 + 4,
    color: colors[Math.floor(Math.random() * colors.length)],
    life: 60,
    maxLife: 60,
  });

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (Math.random() < 0.1) { // 10% chance to create sparkle on mouse move
        const newSparkle = createSparkle(e.clientX, e.clientY);
        setSparkles(prev => [...prev.slice(-20), newSparkle]); // Keep max 20 sparkles
      }
    };

    const handleClick = (e: MouseEvent) => {
      // Create multiple sparkles on click
      const newSparkles = Array.from({ length: 5 }, () => 
        createSparkle(e.clientX, e.clientY)
      );
      setSparkles(prev => [...prev.slice(-15), ...newSparkles]);
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('click', handleClick);

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('click', handleClick);
    };
  }, []);

  useEffect(() => {
    const interval = setInterval(() => {
      setSparkles(prev => 
        prev
          .map(sparkle => ({ ...sparkle, life: sparkle.life - 1 }))
          .filter(sparkle => sparkle.life > 0)
      );
    }, 16); // ~60fps

    return () => clearInterval(interval);
  }, []);

  return (
    <div style={{ position: 'relative', overflow: 'hidden' }}>
      {children}
      
      {sparkles.map(sparkle => (
        <div
          key={sparkle.id}
          style={{
            position: 'fixed',
            left: sparkle.x,
            top: sparkle.y,
            width: sparkle.size,
            height: sparkle.size,
            backgroundColor: sparkle.color,
            borderRadius: '50%',
            pointerEvents: 'none',
            opacity: sparkle.life / sparkle.maxLife,
            transform: `scale(${sparkle.life / sparkle.maxLife})`,
            transition: 'opacity 0.1s ease-out',
            boxShadow: `0 0 ${sparkle.size}px ${sparkle.color}`,
            zIndex: 1000,
          }}
        />
      ))}
    </div>
  );
};