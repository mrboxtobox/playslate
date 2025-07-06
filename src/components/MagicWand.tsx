import React, { useState } from 'react';

interface MagicWandProps {
  onMagicCast: (effect: string) => void;
}

export const MagicWand: React.FC<MagicWandProps> = ({ onMagicCast }) => {
  const [isAnimating, setIsAnimating] = useState(false);

  const magicEffects = [
    { name: 'Rainbow Trail', emoji: '🌈', effect: 'rainbow' },
    { name: 'Sparkle Brush', emoji: '✨', effect: 'sparkle' },
    { name: 'Glow Paint', emoji: '🌟', effect: 'glow' },
    { name: 'Magic Eraser', emoji: '🎭', effect: 'eraser' },
  ];

  const handleMagicClick = (effect: string) => {
    setIsAnimating(true);
    onMagicCast(effect);
    
    setTimeout(() => setIsAnimating(false), 500);
  };

  return (
    <div style={{ 
      display: 'flex', 
      gap: '10px', 
      flexWrap: 'wrap',
      justifyContent: 'center',
      marginBottom: '20px'
    }}>
      {magicEffects.map((magic) => (
        <button
          key={magic.effect}
          onClick={() => handleMagicClick(magic.effect)}
          style={{
            padding: '15px 20px',
            fontSize: '16px',
            backgroundColor: '#9B59B6',
            color: 'white',
            border: 'none',
            borderRadius: '25px',
            cursor: 'pointer',
            fontWeight: 'bold',
            boxShadow: '0 4px 15px rgba(155, 89, 182, 0.3)',
            transform: isAnimating ? 'scale(1.1) rotate(5deg)' : 'scale(1)',
            transition: 'all 0.2s ease',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
          }}
          onMouseEnter={(e) => {
            e.currentTarget.style.transform = 'scale(1.05)';
            e.currentTarget.style.boxShadow = '0 6px 20px rgba(155, 89, 182, 0.4)';
          }}
          onMouseLeave={(e) => {
            if (!isAnimating) {
              e.currentTarget.style.transform = 'scale(1)';
              e.currentTarget.style.boxShadow = '0 4px 15px rgba(155, 89, 182, 0.3)';
            }
          }}
        >
          <span style={{ fontSize: '20px' }}>{magic.emoji}</span>
          {magic.name}
        </button>
      ))}
    </div>
  );
};