import React, { useEffect, useRef, useState } from 'react';
import QrScanner from 'qr-scanner';
import { Camera, CameraOff, RotateCcw } from 'lucide-react';

interface QRScannerProps {
  onScan: (result: string) => void;
  isScanning: boolean;
}

export const QRScanner: React.FC<QRScannerProps> = ({ onScan, isScanning }) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const qrScannerRef = useRef<QrScanner | null>(null);
  const [hasCamera, setHasCamera] = useState(true);
  const [cameraError, setCameraError] = useState<string>('');

  useEffect(() => {
    if (!videoRef.current) return;

    const qrScanner = new QrScanner(
      videoRef.current,
      (result) => {
        onScan(result.data);
      },
      {
        highlightScanRegion: true,
        highlightCodeOutline: true,
        preferredCamera: 'environment',
      }
    );

    qrScannerRef.current = qrScanner;

    if (isScanning) {
      startScanning();
    }

    return () => {
      qrScanner.destroy();
    };
  }, [onScan, isScanning]);

  const startScanning = async () => {
    if (!qrScannerRef.current) return;

    try {
      await qrScannerRef.current.start();
      setCameraError('');
      setHasCamera(true);
    } catch (error) {
      console.error('Camera start error:', error);
      setHasCamera(false);
      setCameraError('Unable to access camera. Please ensure camera permissions are granted.');
    }
  };

  const stopScanning = () => {
    if (qrScannerRef.current) {
      qrScannerRef.current.stop();
    }
  };

  const toggleScanning = () => {
    if (isScanning) {
      stopScanning();
    } else {
      startScanning();
    }
  };

  useEffect(() => {
    if (isScanning) {
      startScanning();
    } else {
      stopScanning();
    }
  }, [isScanning]);

  return (
    <div className="relative">
      <div className="relative overflow-hidden rounded-2xl bg-black shadow-2xl">
        <video
          ref={videoRef}
          className="w-full h-80 object-cover"
          playsInline
          muted
        />
        
        {/* Scanning overlay */}
        <div className="absolute inset-0 border-2 border-blue-500 rounded-2xl">
          <div className="absolute inset-4 border border-white/30 rounded-xl">
            <div className="absolute top-0 left-0 w-6 h-6 border-t-2 border-l-2 border-blue-400"></div>
            <div className="absolute top-0 right-0 w-6 h-6 border-t-2 border-r-2 border-blue-400"></div>
            <div className="absolute bottom-0 left-0 w-6 h-6 border-b-2 border-l-2 border-blue-400"></div>
            <div className="absolute bottom-0 right-0 w-6 h-6 border-b-2 border-r-2 border-blue-400"></div>
          </div>
        </div>

        {/* Status overlay */}
        <div className="absolute top-4 left-4 right-4">
          <div className="bg-black/60 backdrop-blur-sm rounded-lg px-3 py-2">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-2">
                {isScanning ? (
                  <>
                    <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
                    <span className="text-white text-sm font-medium">Scanning...</span>
                  </>
                ) : (
                  <>
                    <div className="w-2 h-2 bg-gray-400 rounded-full"></div>
                    <span className="text-white text-sm font-medium">Paused</span>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Error message */}
      {cameraError && (
        <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <div className="flex items-center space-x-2">
            <CameraOff className="w-5 h-5 text-red-500" />
            <p className="text-red-700 text-sm">{cameraError}</p>
          </div>
        </div>
      )}

      {/* Instructions */}
      <div className="mt-4 text-center">
        <p className="text-gray-600 text-sm">
          Position the QR code within the frame to scan equipment information
        </p>
      </div>
    </div>
  );
};