import { IssuesList } from '../../components/common/IssuesList';

export function AdminIssuesPage() {
  return <IssuesList isAdmin={true} />;
}