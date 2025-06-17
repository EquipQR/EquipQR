export interface Equipment {
  ID: string;
  BusinessID: string;
  Status: string;
  Type: string;
  Location: string;
  MoreFields: {
    manufacturer?: string;
    warranty_expiry?: string;
    [key: string]: any;
  };
  Business: {
    ID: string;
    BusinessName: string;
  };
}