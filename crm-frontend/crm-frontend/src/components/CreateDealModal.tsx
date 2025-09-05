// Replace the contents of CreateDealModal.tsx

import { Modal, Button, Select, Stack, NumberInput, Textarea } from '@mantine/core';
import { useForm, isNotEmpty } from '@mantine/form';
import { useState, useEffect } from 'react';
import { createDeal, updateDeal, getLeads, getProperties } from '../services/apiService';
import type { Deal, CreateDealPayload, Lead, Property } from '../services/apiService';

interface CreateDealModalProps {
  opened: boolean;
  onClose: () => void;
  onSuccess: () => void;
  dealToEdit: Deal | null; // <-- NEW
}

const CreateDealModal = ({ opened, onClose, onSuccess, dealToEdit }: CreateDealModalProps) => {
  const isEditing = !!dealToEdit;
  const [leads, setLeads] = useState<{ value: string; label: string }[]>([]);
  const [properties, setProperties] = useState<{ value: string; label: string }[]>([]);
  const [loadingData, setLoadingData] = useState(true);

  const dealStages = [
    { value: '1', label: 'Prospecting' }, { value: '2', label: 'Qualification' },
    { value: '3', label: 'Negotiation' }, { value: '4' , label: 'Closing' },
  ];

  const form = useForm({
    initialValues: {
      lead_id: '',
      property_id: '',
      stage_id: '1',
      deal_status: 'Pending',
      deal_amount: 0,
      notes: '',
    },
    validate: {
      lead_id: isNotEmpty('A lead must be selected'),
      property_id: isNotEmpty('A property must be selected'),
      deal_amount: (value) => (value > 0 ? null : 'Deal amount must be greater than zero'),
    },
  });

  useEffect(() => {
    if (opened) {
      const fetchData = async () => { /* ... (fetchData is the same as before) ... */ };
      fetchData();

      if (isEditing) {
        form.setValues({
          lead_id: String(dealToEdit.lead_id),
          property_id: String(dealToEdit.property_id),
          stage_id: String(dealToEdit.stage_id),
          deal_status: dealToEdit.deal_status,
          deal_amount: dealToEdit.deal_amount,
          notes: dealToEdit.notes || '',
        });
      } else {
        form.reset();
      }
    }
  }, [dealToEdit, opened]);

  // Fetch data for dropdowns (same as before)
  useEffect(() => {
    if (opened) {
      const fetchData = async () => {
        setLoadingData(true);
        try {
          const [leadsData, propertiesData] = await Promise.all([ getLeads(), getProperties() ]);
          setLeads(leadsData.map((l: Lead) => ({ value: String(l.id), label: `Lead #${l.id} (Contact: ${l.contact_id})` })));
          setProperties(propertiesData.map((p: Property) => ({ value: String(p.id), label: p.name })));
        } catch (error) { console.error("Failed to fetch data for deal form", error); }
        finally { setLoadingData(false); }
      };
      fetchData();
    }
  }, [opened]);


  const handleSubmit = async (values: typeof form.values) => {
    try {
      const payload: CreateDealPayload = {
        lead_id: Number(values.lead_id),
        property_id: Number(values.property_id),
        stage_id: Number(values.stage_id),
        deal_status: values.deal_status,
        deal_amount: values.deal_amount,
        notes: values.notes,
      };
      
      if (isEditing) {
        await updateDeal(dealToEdit.id, payload);
      } else {
        await createDeal(payload);
      }
      onSuccess();
      onClose();
    } catch (error) {
      console.error("Failed to save deal", error);
    }
  };

  return (
    <Modal opened={opened} onClose={onClose} title={isEditing ? 'Edit Deal' : 'Create New Deal'} centered>
      <form onSubmit={form.onSubmit(handleSubmit)}>
        <Stack>
          <Select label="Lead" data={leads} searchable required disabled={loadingData} {...form.getInputProps('lead_id')} />
          <Select label="Property" data={properties} searchable required disabled={loadingData} {...form.getInputProps('property_id')} />
          <Select label="Deal Status" data={['Pending', 'Closed-Won', 'Closed-Lost']} required {...form.getInputProps('deal_status')} />
          <Select label="Deal Stage" data={dealStages} required {...form.getInputProps('stage_id')} />
          <NumberInput label="Deal Amount" prefix="$ " required thousandSeparator min={0} {...form.getInputProps('deal_amount')} />
          <Textarea label="Notes" {...form.getInputProps('notes')} />
          <Button type="submit" mt="md" loading={loadingData}>{isEditing ? 'Update Deal' : 'Create Deal'}</Button>
        </Stack>
      </form>
    </Modal>
  );
};

export default CreateDealModal;