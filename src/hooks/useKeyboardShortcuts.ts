import { useEffect } from 'react';

interface KeyboardShortcutsProps {
  colors: string[];
  brushSize: number;
  setBrushColor: (color: string) => void;
  setBrushSize: (size: number) => void;
  onClearCanvas: () => void;
  onMagicCast: (effect: string) => void;
}

export const useKeyboardShortcuts = ({
  colors,
  brushSize,
  setBrushColor,
  setBrushSize,
  onClearCanvas,
  onMagicCast,
}: KeyboardShortcutsProps) => {
  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      // Prevent default behavior for our shortcuts
      const { key, ctrlKey, metaKey, shiftKey } = event;
      
      // Color shortcuts (1-8)
      if (!isNaN(Number(key)) && Number(key) >= 1 && Number(key) <= 8) {
        const colorIndex = Number(key) - 1;
        if (colors[colorIndex]) {
          setBrushColor(colors[colorIndex]);
          event.preventDefault();
        }
        return;
      }

      // Brush size shortcuts
      if (key === '=' || key === '+') {
        setBrushSize(Math.min(50, brushSize + 5));
        event.preventDefault();
        return;
      }
      
      if (key === '-' || key === '_') {
        setBrushSize(Math.max(2, brushSize - 5));
        event.preventDefault();
        return;
      }

      // Clear canvas (Ctrl/Cmd + C or Delete/Backspace)
      if ((ctrlKey || metaKey) && key.toLowerCase() === 'c') {
        onClearCanvas();
        event.preventDefault();
        return;
      }
      
      if (key === 'Delete' || key === 'Backspace') {
        onClearCanvas();
        event.preventDefault();
        return;
      }

      // Magic shortcuts
      if (shiftKey) {
        switch (key.toLowerCase()) {
          case 'r':
            onMagicCast('rainbow');
            event.preventDefault();
            break;
          case 's':
            onMagicCast('sparkle');
            event.preventDefault();
            break;
          case 'g':
            onMagicCast('glow');
            event.preventDefault();
            break;
          case 'e':
            onMagicCast('eraser');
            event.preventDefault();
            break;
        }
        return;
      }

      // Quick brush size presets
      switch (key.toLowerCase()) {
        case 't': // Tiny
          setBrushSize(4);
          event.preventDefault();
          break;
        case 's': // Small
          setBrushSize(8);
          event.preventDefault();
          break;
        case 'm': // Medium
          setBrushSize(16);
          event.preventDefault();
          break;
        case 'l': // Large
          setBrushSize(32);
          event.preventDefault();
          break;
        case 'x': // eXtra large
          setBrushSize(50);
          event.preventDefault();
          break;
      }
    };

    // Add event listener
    document.addEventListener('keydown', handleKeyPress);

    // Cleanup
    return () => {
      document.removeEventListener('keydown', handleKeyPress);
    };
  }, [colors, brushSize, setBrushColor, setBrushSize, onClearCanvas, onMagicCast]);
};