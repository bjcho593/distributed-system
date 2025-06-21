import express from 'express';
import dotenv from 'dotenv';
import authRoutes from './routes/auth.js';

dotenv.config();
const app = express();

app.use(express.json());
app.use('/auth', authRoutes);

app.get('/auth/status', (req, res) => {
  res.json({ message: 'Auth Service is running ✅' });
});

const PORT = 8082;
app.listen(PORT, () => {
  console.log(`Auth Service running on port ${PORT}`);
});
