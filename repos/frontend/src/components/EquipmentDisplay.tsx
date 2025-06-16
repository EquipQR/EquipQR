import React from 'react';
import { Equipment } from '../types/equipment';
import { 
  Building2, 
  MapPin, 
  Wrench, 
  Calendar, 
  Shield, 
  Info,
  CheckCircle,
  AlertCircle,
  Clock
} from 'lucide-react';

interface EquipmentDisplayProps {
  equipment: Equipment | null;
}

export const EquipmentDisplay: React.FC<EquipmentDisplayProps> = ({ equipment }) => {
  if (!equipment) {
    return (
      <div className="bg-white rounded-2xl shadow-lg p-8 text-center">
        <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
          <Wrench className="w-8 h-8 text-gray-400" />
        </div>
        <h3 className="text-lg font-semibold text-gray-900 mb-2">No Equipment Scanned</h3>
        <p className="text-gray-600">Scan a QR code to view equipment details</p>
      </div>
    );
  }

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'in service':
        return 'bg-green-100 text-green-800 border-green-200';
      case 'out of service':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'maintenance':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'in service':
        return <CheckCircle className="w-4 h-4" />;
      case 'out of service':
        return <AlertCircle className="w-4 h-4" />;
      case 'maintenance':
        return <Clock className="w-4 h-4" />;
      default:
        return <Info className="w-4 h-4" />;
    }
  };

  const formatDate = (dateString: string) => {
    try {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      });
    } catch {
      return dateString;
    }
  };

  return (
    <div className="space-y-6">
      {/* Equipment Header */}
      <div className="bg-white rounded-2xl shadow-lg p-6">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h3 className="text-xl font-bold text-gray-900 mb-1">
              {equipment.Type.charAt(0).toUpperCase() + equipment.Type.slice(1)}
            </h3>
            <p className="text-sm text-gray-600">ID: {equipment.ID}</p>
          </div>
          <div className={`flex items-center space-x-2 px-3 py-2 rounded-full border ${getStatusColor(equipment.Status)}`}>
            {getStatusIcon(equipment.Status)}
            <span className="text-sm font-medium capitalize">{equipment.Status}</span>
          </div>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
              <Building2 className="w-5 h-5 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Business</p>
              <p className="font-medium text-gray-900">{equipment.Business.BusinessName}</p>
            </div>
          </div>
          
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
              <MapPin className="w-5 h-5 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-600">Location</p>
              <p className="font-medium text-gray-900">{equipment.Location}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Equipment Details */}
      <div className="bg-white rounded-2xl shadow-lg p-6">
        <h4 className="text-lg font-semibold text-gray-900 mb-4">Equipment Details</h4>
        
        <div className="space-y-4">
          {equipment.MoreFields.manufacturer && (
            <div className="flex items-center justify-between py-3 border-b border-gray-100">
              <div className="flex items-center space-x-2">
                <Wrench className="w-4 h-4 text-gray-400" />
                <span className="text-gray-600">Manufacturer</span>
              </div>
              <span className="font-medium text-gray-900">{equipment.MoreFields.manufacturer}</span>
            </div>
          )}
          
          {equipment.MoreFields.warranty_expiry && (
            <div className="flex items-center justify-between py-3 border-b border-gray-100">
              <div className="flex items-center space-x-2">
                <Shield className="w-4 h-4 text-gray-400" />
                <span className="text-gray-600">Warranty Expiry</span>
              </div>
              <span className="font-medium text-gray-900">
                {formatDate(equipment.MoreFields.warranty_expiry)}
              </span>
            </div>
          )}
          
          <div className="flex items-center justify-between py-3 border-b border-gray-100">
            <div className="flex items-center space-x-2">
              <Calendar className="w-4 h-4 text-gray-400" />
              <span className="text-gray-600">Equipment Type</span>
            </div>
            <span className="font-medium text-gray-900 capitalize">{equipment.Type}</span>
          </div>
          
          <div className="flex items-center justify-between py-3">
            <div className="flex items-center space-x-2">
              <Building2 className="w-4 h-4 text-gray-400" />
              <span className="text-gray-600">Business ID</span>
            </div>
            <span className="font-medium text-gray-900 text-sm">{equipment.BusinessID}</span>
          </div>
        </div>
      </div>

      {/* Additional Fields */}
      {Object.keys(equipment.MoreFields).length > 2 && (
        <div className="bg-white rounded-2xl shadow-lg p-6">
          <h4 className="text-lg font-semibold text-gray-900 mb-4">Additional Information</h4>
          <div className="space-y-3">
            {Object.entries(equipment.MoreFields).map(([key, value]) => {
              if (key === 'manufacturer' || key === 'warranty_expiry') return null;
              
              return (
                <div key={key} className="flex items-center justify-between py-2 border-b border-gray-100">
                  <span className="text-gray-600 capitalize">{key.replace(/_/g, ' ')}</span>
                  <span className="font-medium text-gray-900">{value}</span>
                </div>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
};