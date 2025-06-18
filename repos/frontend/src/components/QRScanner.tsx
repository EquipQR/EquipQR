import React, { useEffect, useRef, useState } from 'react';
import QrScanner from 'qr-scanner';
import { Flashlight, FlashlightOff, Info, Loader2 } from 'lucide-react';

interface QRScannerProps {
  onScan: (result: string) => void;
  isScanning: boolean;
  loading?: boolean;
  error?: string;
}

export const QRScanner: React.FC<QRScannerProps> = ({ onScan, loading, error }) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const qrScannerRef = useRef<QrScanner | null>(null);
  const torchTrackRef = useRef<MediaStreamTrack | null>(null);
  const [_, setHasCamera] = useState(true);
  const [cameraError, setCameraError] = useState<string>('');
  const [torchSupported, setTorchSupported] = useState(false);
  const [torchEnabled, setTorchEnabled] = useState(false);
  const [torchError, setTorchError] = useState<string | null>(null);

  useEffect(() => {
    const setup = async () => {
      try {
        const devices = await navigator.mediaDevices.enumerateDevices();
        const videoDevices = devices.filter((d) => d.kind === 'videoinput');
        if (videoDevices.length === 0) throw new Error('No video input devices found.');

        const selectedDevice = videoDevices[videoDevices.length - 1];

        const stream = await navigator.mediaDevices.getUserMedia({
          video: { deviceId: { exact: selectedDevice.deviceId } }
        });

        const track = stream.getVideoTracks()[0];
        const capabilities = track.getCapabilities?.() ?? {};
        const supportsTorch = 'torch' in capabilities;

        torchTrackRef.current = track;
        setTorchSupported(supportsTorch);

        if (videoRef.current) {
          videoRef.current.srcObject = stream;
          videoRef.current.onloadedmetadata = () => {
            videoRef.current?.play().catch(console.error);
          };
        }
        const qrScanner = new QrScanner(
          videoRef.current!,
          // @ts-expect-error
          (result) => onScan(result.data),
          {
            highlightScanRegion: true,
            highlightCodeOutline: true,
            deviceId: selectedDevice.deviceId
          }
        );
        qrScannerRef.current = qrScanner;
        await qrScanner.start();

        setHasCamera(true);
        setCameraError('');
      } catch (err) {
        console.error(err);
        setHasCamera(false);
        setCameraError('Unable to access camera. Please ensure permissions are granted.');
      }
    };

    setup();

    return () => {
      qrScannerRef.current?.destroy();
      torchTrackRef.current?.stop();
      const stream = videoRef.current?.srcObject as MediaStream | null;
      stream?.getTracks().forEach((t) => t.stop());
    };
  }, [onScan]);

  const toggleTorch = async () => {
    if (!torchTrackRef.current) {
      setTorchError('No torch-capable track available.');
      return;
    }

    try {
      const enableTorch = !torchEnabled;
      await torchTrackRef.current?.applyConstraints({
        advanced: [{ torch: enableTorch } as any]
      });
      setTorchEnabled(enableTorch);
      setTorchError(null);
    } catch (err) {
      console.error(err);
      setTorchError('Failed to toggle torch.');
    }

  };

  return (
    <div className="fixed inset-0 bg-black">
      <div className="relative w-full h-full">
        <video
          ref={videoRef}
          className="w-full h-full object-cover"
          playsInline
          muted
        />

        <div className="absolute top-0 left-0 right-0 p-6 z-10 flex items-center justify-between">
          <div className="bg-black/60 backdrop-blur-sm rounded-full px-4 py-2">
            <div className="flex items-center space-x-2">
              <div className="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
              <span className="text-white text-sm font-medium">Scanning for QR codes...</span>
            </div>
          </div>

          {torchSupported && (
            <button
              onClick={toggleTorch}
              className={`p-3 rounded-full backdrop-blur-sm transition-all ${torchEnabled
                  ? 'bg-yellow-500/80 text-white'
                  : 'bg-black/60 text-white hover:bg-black/80'
                }`}
            >
              {torchEnabled ? (
                <Flashlight className="w-6 h-6" />
              ) : (
                <FlashlightOff className="w-6 h-6" />
              )}
            </button>
          )}
        </div>

        <div className="absolute bottom-0 left-0 right-0 p-6 z-10">
          <div className="text-center">
            <div className="bg-black/60 backdrop-blur-sm rounded-2xl px-6 py-4 inline-block">
              <h2 className="text-white text-xl font-semibold mb-2">Scan Equipment QR Code</h2>
              <p className="text-white/80 text-sm">
                Position the QR code within the frame to scan equipment information
              </p>
            </div>
          </div>
        </div>

        {loading && (
          <div className="absolute inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-20">
            <div className="bg-white rounded-2xl p-8 text-center max-w-sm mx-4">
              <Loader2 className="w-12 h-12 text-blue-600 animate-spin mx-auto mb-4" />
              <h3 className="text-lg font-semibold text-gray-900 mb-2">Loading Equipment</h3>
              <p className="text-gray-600">Fetching equipment information...</p>
            </div>
          </div>
        )}

        {(error || cameraError) && (
          <div className="absolute inset-0 bg-black/80 backdrop-blur-sm flex items-center justify-center z-20">
            <div className="bg-white rounded-2xl p-8 text-center max-w-sm mx-4">
              <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Info className="w-8 h-8 text-red-600" />
              </div>
              <h3 className="text-lg font-semibold text-gray-900 mb-2">
                {cameraError ? 'Camera Error' : 'Equipment Not Found'}
              </h3>
              <p className="text-gray-600 mb-6">
                {cameraError || error}
              </p>
              <button
                onClick={() => window.location.reload()}
                className="w-full bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors font-medium"
              >
                Try Again
              </button>
            </div>
          </div>
        )}

        {torchError && (
          <div className="absolute top-20 left-0 right-0 px-6 z-30 flex justify-center">
            <div className="bg-red-600 text-white text-sm rounded-lg px-4 py-2 shadow-lg">
              {torchError}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
