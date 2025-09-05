// File: src/components/RecentActivity.tsx
import { Timeline, Text } from '@mantine/core';
import { IconGitBranch, IconGitPullRequest, IconGitCommit, IconMessageDots } from '@tabler/icons-react';
import { useData } from '../context/DataContext';

const RecentActivity = () => {
  // We get the raw data from our global context
  const { leads, deals, loading } = useData();

  if (loading) return null;

  // Combine and sort the latest 5 activities
  const allActivities = [
    ...leads.map(l => ({ type: 'New Lead', date: new Date(l.created_at), id: `l-${l.id}` })),
    ...deals.map(d => ({ type: 'New Deal', date: new Date(d.created_at), id: `d-${d.id}` })),
  ];

  const sortedActivities = allActivities
    .sort((a, b) => b.date.getTime() - a.date.getTime())
    .slice(0, 5); // Get the 5 most recent items

  return (
    <Timeline active={sortedActivities.length} bulletSize={24} lineWidth={2}>
      {sortedActivities.map(activity => (
        <Timeline.Item 
            bullet={activity.type === 'New Lead' ? <IconGitPullRequest size={12} /> : <IconGitCommit size={12} />} 
            title={activity.type}
            key={activity.id}
        >
          <Text c="dimmed" size="xs">
            {activity.date.toLocaleString()}
          </Text>
        </Timeline.Item>
      ))}
    </Timeline>
  );
};

export default RecentActivity;