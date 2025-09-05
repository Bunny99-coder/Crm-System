// Replace the entire contents of src/context/DataContext.tsx
import React, { createContext, useState, useContext, useEffect, useMemo } from 'react';
import type { ReactNode } from 'react';
import { getUsers, getContacts, getProperties, getLeads, getDeals } from '../services/apiService';
import type { UserSelectItem, Contact, Property, Lead, Deal } from '../services/apiService';
import { useAuth } from './AuthContext';

interface DataContextType {
  userMap: Map<number, string>;
  contactMap: Map<number, string>;
  propertyMap: Map<number, string>;
  leadMap: Map<number, string>;
  deals: Deal[]; // RecentActivity needs the full deal list
  leads: Lead[]; // RecentActivity needs the full lead list
  loading: boolean;
}

const DataContext = createContext<DataContextType | undefined>(undefined);

export const DataProvider = ({ children }: { children: ReactNode }) => {
  const { isAuthenticated } = useAuth();
  const [loading, setLoading] = useState(false);
  
  const [userMap, setUserMap] = useState(new Map<number, string>());
  const [contactMap, setContactMap] = useState(new Map<number, string>());
  const [propertyMap, setPropertyMap] = useState(new Map<number, string>());
  const [leadMap, setLeadMap] = useState(new Map<number, string>());
  const [deals, setDeals] = useState<Deal[]>([]);
  const [leads, setLeads] = useState<Lead[]>([]);

  useEffect(() => {
    if (isAuthenticated) {
      setLoading(true);
      Promise.all([
        getUsers(),
        getContacts(),
        getProperties(),
        getLeads(),
        getDeals(),
      ])
        .then(([usersData, contactsData, propertiesData, leadsData, dealsData]) => {
          const safeContacts = contactsData || [];
          const newContactMap = new Map(safeContacts.map(c => [c.id, `${c.first_name} ${c.last_name}`]));
          
          const safeUsers = usersData || [];
          setUserMap(new Map(safeUsers.map(u => [u.id, u.username])));
          setContactMap(newContactMap);

          const safeProperties = propertiesData || [];
          setPropertyMap(new Map(safeProperties.map(p => [p.id, p.name])));

          const safeLeads = leadsData || [];
          setLeads(safeLeads);
          setLeadMap(new Map(safeLeads.map(l => [l.id, newContactMap.get(l.contact_id) || `Lead #${l.id}`])));

          const safeDeals = dealsData || [];
          setDeals(safeDeals);
        })
        .catch(error => console.error("Failed to load global data for context", error))
        .finally(() => setLoading(false));
    } else {
      // Clear data on logout
      setUserMap(new Map());
      setContactMap(new Map());
      setPropertyMap(new Map());
      setLeadMap(new Map());
      setDeals([]);
      setLeads([]);
    }
  }, [isAuthenticated]);

  // THIS useMemo HOOK IS THE FIX.
  // It ensures the 'value' object is only recreated when the data actually changes,
  // preventing an infinite re-render loop.
  const value = useMemo(() => ({
    userMap, contactMap, propertyMap, leadMap, deals, leads, loading
  }), [userMap, contactMap, propertyMap, leadMap, deals, leads, loading]);

  return (
    <DataContext.Provider value={value}>
      {children}
    </DataContext.Provider>
  );
};

export const useData = () => {
  const context = useContext(DataContext);
  if (context === undefined) {
    throw new Error('useData must be used within a DataProvider');
  }
  return context;
};