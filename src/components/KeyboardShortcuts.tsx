import React, { useState } from 'react';

export const KeyboardShortcuts: React.FC = () => {
  const [isVisible, setIsVisible] = useState(false);

  const shortcuts = [
    { keys: '1-8', description: 'Select colors' },
    { keys: '+/-', description: 'Increase/decrease brush size' },
    { keys: 'T', description: 'Tiny brush (4px)' },
    { keys: 'S', description: 'Small brush (8px)' },
    { keys: 'M', description: 'Medium brush (16px)' },
    { keys: 'L', description: 'Large brush (32px)' },
    { keys: 'X', description: 'Extra large brush (50px)' },
    { keys: 'Ctrl/⌘ + C', description: 'Clear canvas' },
    { keys: 'Delete', description: 'Clear canvas' },
    { keys: 'Shift + R', description: 'Rainbow mode' },
    { keys: 'Shift + S', description: 'Sparkle mode' },
    { keys: 'Shift + G', description: 'Glow mode' },
    { keys: 'Shift + E', description: 'Eraser mode' },
  ];

  return (
    <div style={{ position: 'fixed', top: '20px', right: '20px', zIndex: 1000 }}>
      <button
        onClick={() => setIsVisible(!isVisible)}
        style={{
          padding: '10px 15px',
          backgroundColor: '#6A4C93',
          color: 'white',
          border: 'none',
          borderRadius: '20px',
          cursor: 'pointer',
          fontWeight: '600',
          fontSize: '14px',
          boxShadow: '0 4px 12px rgba(106, 76, 147, 0.3)',
          transition: 'all 0.2s ease',
        }}
        onMouseEnter={(e) => {
          e.currentTarget.style.transform = 'scale(1.05)';
          e.currentTarget.style.boxShadow = '0 6px 16px rgba(106, 76, 147, 0.4)';
        }}
        onMouseLeave={(e) => {
          e.currentTarget.style.transform = 'scale(1)';
          e.currentTarget.style.boxShadow = '0 4px 12px rgba(106, 76, 147, 0.3)';
        }}
      >
        ⌨️ Shortcuts
      </button>

      {isVisible && (
        <div
          style={{
            position: 'absolute',
            top: '60px',
            right: '0',
            backgroundColor: 'white',
            borderRadius: '15px',
            padding: '20px',
            minWidth: '280px',
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.15)',
            border: '2px solid #E5E5E5',
          }}
        >
          <h3 
            style={{ 
              margin: '0 0 15px 0', 
              color: '#6A4C93',
              fontSize: '18px',
              fontWeight: '700',
              textAlign: 'center'
            }}
          >
            Keyboard Shortcuts
          </h3>
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
            {shortcuts.map((shortcut, index) => (
              <div 
                key={index}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '8px 12px',
                  backgroundColor: '#F8F9FA',
                  borderRadius: '8px',
                  fontSize: '14px',
                }}
              >
                <span style={{ fontWeight: '600', color: '#495057' }}>
                  {shortcut.keys}
                </span>
                <span style={{ color: '#6C757D' }}>
                  {shortcut.description}
                </span>
              </div>
            ))}
          </div>
          
          <div 
            style={{
              marginTop: '15px',
              padding: '10px',
              backgroundColor: '#E3F2FD',
              borderRadius: '8px',
              fontSize: '12px',
              color: '#1976D2',
              textAlign: 'center',
            }}
          >
            💡 Tip: Use number keys 1-8 to quickly switch colors!
          </div>
        </div>
      )}
    </div>
  );
};