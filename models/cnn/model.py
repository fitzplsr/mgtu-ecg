import os
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.data import Dataset, DataLoader
import pyedflib
import numpy as np
from sklearn.model_selection import train_test_split
from sklearn.metrics import accuracy_score, roc_auc_score
from imblearn.over_sampling import SMOTE
from tqdm import tqdm
import matplotlib.pyplot as plt

# Параметры
EDF_DIR = "/content/drive/MyDrive/PulmHypert/PulmHypert"
SEQ_LENGTHS = [2500]
BATCH_SIZE = 16
NUM_WORKERS = 2
LR = 1e-3
WEIGHT_DECAY = 1e-5
EPOCHS = 50
EARLY_STOP_PATIENCE = 5
DEVICE = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

# Focal Loss
class FocalLoss(nn.Module):
    def __init__(self, alpha=1, gamma=2, reduction='mean'):
        super().__init__()
        self.alpha, self.gamma, self.reduction = alpha, gamma, reduction
    def forward(self, inputs, targets):
        bce = nn.functional.binary_cross_entropy_with_logits(inputs, targets, reduction='none')
        probs = torch.sigmoid(inputs)
        pt = torch.where(targets==1, probs, 1-probs)
        loss = self.alpha * (1-pt)**self.gamma * bce
        return loss.mean() if self.reduction=='mean' else loss.sum()

# 1D Residual Block
class ResBlock1D(nn.Module):
    def __init__(self, in_ch, out_ch, k=3, pad=1):
        super().__init__()
        self.conv1 = nn.Conv1d(in_ch, out_ch, k, padding=pad)
        self.bn1 = nn.BatchNorm1d(out_ch)
        self.conv2 = nn.Conv1d(out_ch, out_ch, k, padding=pad)
        self.bn2 = nn.BatchNorm1d(out_ch)
        self.relu = nn.ReLU(inplace=True)
        self.skip = nn.Conv1d(in_ch, out_ch, 1) if in_ch!=out_ch else nn.Identity()
    def forward(self, x):
        res = self.skip(x)
        out = self.relu(self.bn1(self.conv1(x)))
        out = self.bn2(self.conv2(out))
        out += res
        return self.relu(out)

# Advanced Model: ResNet + LSTM
class AdvancedECGModel(nn.Module):
    def __init__(self, seq_length):
        super().__init__()
        self.resnet = nn.Sequential(
            ResBlock1D(1, 16), nn.MaxPool1d(2),
            ResBlock1D(16, 32), nn.MaxPool1d(2),
            ResBlock1D(32, 64), nn.MaxPool1d(2)
        )
        red_len = seq_length // 8
        self.lstm = nn.LSTM(64, 64, batch_first=True, bidirectional=True)
        self.classifier = nn.Sequential(
            nn.Linear(2*64*red_len, 128), nn.ReLU(inplace=True),
            nn.Dropout(0.5), nn.Linear(128, 1)
        )
    def forward(self, x):
        out = self.resnet(x)
        out = out.permute(0, 2, 1)
        out, _ = self.lstm(out)
        out = out.reshape(out.size(0), -1)
        return self.classifier(out).squeeze(1)

# Dataset with Augmentation
class ECGDataset(Dataset):
    def __init__(self, signals, labels=None, augment=False):
        self.signals = signals
        self.labels = labels
        self.augment = augment
    def __len__(self): return len(self.signals)
    def __getitem__(self, idx):
        sig = self.signals[idx].astype(np.float32)
        # augmentations for training
        if self.augment:
            # random time shift
            shift = np.random.randint(-100, 100)
            sig = np.roll(sig, shift)
            # random scaling
            sig *= np.random.uniform(0.9, 1.1)
            # gaussian noise
            sig += np.random.normal(0, 0.01, size=sig.shape)
        # normalization
        sig = (sig - sig.mean()) / (sig.std() + 1e-6)
        sig = torch.from_numpy(sig).unsqueeze(0)
        if self.labels is not None:
            return sig, torch.tensor(self.labels[idx], dtype=torch.float32)
        return sig

# Load signal
def load_signal(path, seq_length):
    with pyedflib.EdfReader(path) as f:
        ch = f.getSignalLabels();
        idx = ch.index('ECG V1')
        sig = f.readSignal(idx)
    return sig[:seq_length]

# Compute metrics
def compute_metrics(model, loader):
    model.eval()
    probs, targets = [], []
    with torch.no_grad():
        for X, y in loader:
            X = X.to(DEVICE)
            out = model(X)
            probs.extend(torch.sigmoid(out).cpu().numpy())
            targets.extend(y.numpy())
    preds = [p>0.5 for p in probs]
    return accuracy_score(targets, preds), roc_auc_score(targets, probs)
