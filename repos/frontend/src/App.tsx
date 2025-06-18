import { useState } from 'react';
import { QRScanner } from './components/QRScanner';
import { EquipmentDisplay } from './components/EquipmentDisplay';
import { Equipment } from './types/equipment';
import { fetchEquipment } from './services/equipmentApi';
import { Settings, Scan, Info } from 'lucide-react';

function App() {
  const [isScanning, setIsScanning] = useState(false);
  const [equipment, setEquipment] = useState<Equipment | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string>('');
  const handleScan = async (scannedId: string) => {
    setLoading(true);
    setError('');

    try {
      const equipmentData = await fetchEquipment(scannedId);
      setEquipment(equipmentData);
      setIsScanning(false);
    } catch (err) {
      if (err instanceof TypeError && err.message.includes('Failed to fetch')) {
        setError('Network error or CORS issue. Check the server and CORS configuration.');
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Unexpected error occurred while fetching equipment data.');
      }
      setEquipment(null);
    } finally {
      setLoading(false);
    }
  };


  const startNewScan = () => {
    setEquipment(null);
    setError('');
    setIsScanning(true);
  };

  const toggleScanning = () => {
    setIsScanning(!isScanning);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="bg-blue-600 p-2 rounded-lg">
                <Scan className="w-6 h-6 text-white" />
              </div>
              <div>
                <h1 className="text-xl font-bold text-gray-900">Equipment Scanner</h1>
                <p className="text-sm text-gray-600">Scan QR codes to view equipment details</p>
              </div>
            </div>
            <button className="p-2 text-gray-400 hover:text-gray-600 transition-colors">
              <Settings className="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Scanner Section */}
          <div className="space-y-6">
            <div className="bg-white rounded-2xl shadow-lg p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-semibold text-gray-900">QR Code Scanner</h2>
                <button
                  onClick={toggleScanning}
                  className={`px-4 py-2 rounded-lg font-medium transition-all ${isScanning
                      ? 'bg-red-100 text-red-700 hover:bg-red-200'
                      : 'bg-blue-100 text-blue-700 hover:bg-blue-200'
                    }`}
                >
                  {isScanning ? 'Stop Scanning' : 'Start Scanning'}
                </button>
              </div>

              <QRScanner onScan={handleScan} isScanning={isScanning} />

              {loading && (
                <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
                    <p className="text-blue-700 font-medium">Loading equipment data...</p>
                  </div>
                </div>
              )}

              {error && (
                <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
                  <div className="flex items-center space-x-2">
                    <Info className="w-5 h-5 text-red-500 flex-shrink-0" />
                    <p className="text-red-700 text-sm">{error}</p>
                  </div>
                </div>
              )}
            </div>

            {/* Quick Actions */}
            <div className="bg-white rounded-2xl shadow-lg p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
              <div className="space-y-3">
                <button
                  onClick={startNewScan}
                  className="w-full flex items-center justify-center space-x-2 bg-blue-600 text-white py-3 px-4 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  <Scan className="w-4 h-4" />
                  <span>New Scan</span>
                </button>
                <button
                  onClick={() => setEquipment(null)}
                  className="w-full flex items-center justify-center space-x-2 bg-gray-100 text-gray-700 py-3 px-4 rounded-lg hover:bg-gray-200 transition-colors"
                >
                  <span>Clear Results</span>
                </button>
              </div>
            </div>
          </div>

          {/* Equipment Display Section */}
          <div className="space-y-6">
            <EquipmentDisplay equipment={equipment} />
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;