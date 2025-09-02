---
title: "[LLM] What is Transformer?"
author:
    name: "Jungbeom"
    url: "https://github.com/JungbeomLee"
date: "2025-09-02"
readingTime: "10분"
excerpt: "Transformer는 Attention 메커니즘을 활용하여 RNN/LSTM의 한계를 극복한 딥러닝 모델입니다. Self-Attention, Multi-Head Attention, LayerNorm, Residual Connection 등을 통해 문맥 이해와 병렬 연산을 가능하게 하며, GPT·BERT·T5 등 대형 언어 모델의 기반이 되었습니다."
---

이 글은 **[LLM] 시리즈** 중 하나로, Transformer의 개념을 정리한 글입니다.

## What is Attention?

Attention은 **가중치(weight)를 사용하여 어떤 입력이 더 중요한지 계산하고 그 정보를 반영하는 과정**입니다.  
즉, 주어진 입력 데이터에서 각 요소가 출력에 얼마나 중요한지를 학습하여 가중치를 다르게 부여하는 방식입니다.

## What is Transformer?

Transformer는 **자연어 처리(NLP) 및 시퀀스 데이터 처리**에 사용되는 인공지능(AI) 모델로, Attention 메커니즘을 활용하여 문맥을 이해하고 병렬 연산이 가능하여 높은 성능을 보입니다.  
대형 언어 모델(LLM)들의 기반이 되는 가장 강력한 딥러닝 모델 중 하나입니다.  
- PS |
    - 기존 RNN(Recurrent Neural Network)과 LSTM(Long Short-Term Memory)의 한계를 극복  
    - 기계 번역, 텍스트 요약, 문장 생성 등에서 SOTA(State-of-the-art) 성능 달성  
- 핵심 개념
    - Self-Attention으로 단어 간 관계 학습  
    - 병렬 연산 가능 → 학습 속도 향상  
    - Positional Encoding으로 단어 순서 학습  
    - Encoder-Decoder 구조 사용 (GPT는 Decoder-only)  

#### Transformer 특징

- **Self-Attention Mechanism**: 문장의 각 단어가 다른 단어들과의 관계를 학습  
- **병렬 연산(Parallelization)**: 전체 문장을 동시에 입력받아 빠른 학습 가능  
- **Positional Encoding**: 단어 순서 정보 반영  
- **확장성**: 대규모 데이터에서 강력한 성능  
- GPT, BERT, T5 같은 LLM들이 모두 Transformer 기반  

#### Transformer 구조

| **입력 경로**           | **Encoder (N×)**                  | **Decoder (N×)**                          | **출력 경로**        |
|--------------------------|-----------------------------------|--------------------------------------------|----------------------|
| Inputs                   | Multi-Head Attention (Self)       | Masked Multi-Head Attention (Self)         | Linear → Softmax     |
| Input Embedding          | Add & Norm                        | Add & Norm                                 | Output Probabilities |
| Positional Encoding      | Feed Forward (MLP)                | Multi-Head Attention (Cross w/ Encoder)    |                      |
|                          | Add & Norm                        | Add & Norm + Feed Forward                  |                      |
|                          | (반복 N×)                         | (반복 N×)                                  |                      |
- Encoder
    - 입력 문장을 토큰화(Tokenization) 후 임베딩(Embedding)  
    - Positional Encoding 추가  
    - Self-Attention 수행 → 문맥 학습  
    - Feedforward Neural Network (FFN) 적용  
    - Layer Normalization, Residual Connection 사용  
- Decoder
    - Encoder 출력 벡터 입력  
    - Masked Self-Attention 적용 (이전 단어까지만 참고)  
    - Encoder 출력과 결합한 Attention 수행  
    - Feedforward NN, Softmax로 다음 단어 예측  
- Self-Attention 메커니즘
    - 입력 단어 벡터 X를 Query(Q), Key(K), Value(V)로 변환합니다.
    - Query: "이 단어가 다른 단어와 얼마나 연관이 있는가?"  
    - Key: "이 단어가 어떤 정보를 갖고 있는가?"  
    - Value: "그 정보를 모델이 어떻게 활용할 것인가?"  

## Attention Score 계산
- Query와 Key 내적(dot product)  
- 차원 수 \( d_k \)로 스케일링  
- Softmax로 확률 분포 계산  
- 예시: “sat”은 “cat”과 “on”에 높은 가중치 → 문맥적으로 연관성 큼 
최종 출력은 Value 벡터의 가중합으로 생성됩니다.

## Multi-Head Attention
- 단일 Self-Attention은 다의적 의미 반영에 한계  
- Multi-Head Attention은 여러 Head가 서로 다른 관계를 학습  
- Head별 결과를 Concat 후 가중치 \( W^O \) 적용해 최종 벡터 생성  

## Feedforward Neural Network (FFN)
- 단순 Fully Connected Layer  
- ReLU 활성화 함수 적용  
- Self-Attention 결과 벡터를 비선형 변환해 표현력 강화  

## Layer Normalization
- 입력 벡터를 평균과 표준편차로 정규화  
- BatchNorm과 달리 배치 크기에 무관  
- Gradient Vanishing/Exploding 방지  
- Transformer, RNN 등에 최적화  

### 정규화 수식
\[\hat{x_i} = \frac{x_i - \mu}{\sigma + \epsilon}\]  
- \(\mu\): 입력 평균  
- \(\sigma\): 입력 표준편차  
- \(\gamma, \beta\): 학습 가능한 스케일링/이동 파라미터  

## Residual Connection
\[y = f(x) + x\]  
- 입력을 변환된 출력과 더하는 Skip Connection  
- 기울기 소실 문제 해결  
- 학습 속도 향상 및 정보 손실 방지  

## Transformer 최종 출력
- Decoder 마지막 벡터 → Linear Projection (행렬 \( W_o \))  
- 단어 사전 크기 V에 맞게 Logits 계산  
- Softmax로 확률 분포 변환  
- Ex :

| 단어   | Logits | 확률 |
|--------|--------|------|
| Paris  | 2.5    | 0.7  |
| London | 1.8    | 0.2  |
| Tokyo  | 0.5    | 0.08 |
| Berlin | -0.2   | 0.02 |
가장 확률이 높은 단어가 최종 출력으로 선택되고, 다음 입력으로 추가되어 반복적으로 문장이 완성됩니다.
