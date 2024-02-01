import { useNavigate } from 'react-router-dom';

const NoFoundPage = () => {
  const navigate = useNavigate();
  return (
    <div>
      <button onClick={() => navigate('/')}>Back Home</button>
    </div>
  );
};

export default NoFoundPage;
