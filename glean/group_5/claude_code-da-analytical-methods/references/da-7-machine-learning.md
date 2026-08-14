<!-- hub-reference-banner -->
> **Reference file — part of the `da-analytical-methods` hub.** Formerly the standalone `da-7-machine-learning` skill.
> Sibling topics in this family are now reference files under the hubs (`da-1-foundations-theory`, `da-2-data-analysis-lifecycle`, `da-3-data-acquisition-sampling`, `da-analytical-methods`, `da-data-engineering-platform`, `da-applied-and-communication`) — **not** standalone
> skills. Ignore any "use the X skill" / `related_skills` / SKIP pointers below that name a bare sibling
> skill; load that topic's `references/<name>.md` from the owning hub (see the hub's "Cross-hub map").

---

---
name: da-7-machine-learning
description: Curriculum reference for Machine Learning — section 7 of the data analysis curriculum. Covers the ML taxonomy (supervised, unsupervised, reinforcement, self-supervised), bias-variance tradeoff and regularization, deep learning architectures (CNNs, RNNs, Transformers), foundation models and the 2026 LLM landscape (Claude Opus 4.7, GPT-5.5, Gemini 3.1 Pro, Llama 4, DeepSeek V3.2), hyperparameter tuning (Optuna, Ray Tune, BOHB, Bayesian optimization), evaluation methodology for both classical ML (precision/recall/F1, ROC-AUC, regression metrics) and LLMs (HELM, MMLU, MT-Bench, Chatbot Arena, SWE-bench), and MLOps (drift detection, train-serve skew, MLflow, W&B, feature stores). Use as a study companion when reviewing the curriculum or when a project demands first-principles grounding before applying ML/DL techniques. TRIGGER: questions about ML taxonomy and learning paradigms; bias-variance, overfitting/underfitting, regularization choices; selecting between CNNs, RNNs, and Transformers; foundation model selection in 2026; hyperparameter search strategy or tool choice (Optuna vs Ray Tune vs HyperBand); choosing evaluation metrics for classification, regression, ranking, or LLMs; designing drift monitors, train-serve skew prevention, or experiment tracking; curriculum review covering "section 7" of the data analysis path. SKIP: pure statistical inference with no model fitting (use da-1-4-statistical-inference-foundations); data cleaning, EDA, or feature engineering steps that precede modeling (use da-4-data-cleaning-preparation); end-to-end visualization and reporting after a model is trained (use da-8-data-visualization or da-9-reporting-communication); production prompt engineering with a finished LLM and no model selection question (use prompt-engineering); building a RAG system architecture (use rag-architecture); MongoDB-specific ML/AI features such as Atlas Vector Search or Atlas Search (use mongodb-atlas-vector-search or mongodb-atlas-search); MLOps tied to a specific cloud platform with no model-selection question (use the relevant cloud skill).
when_to_use: Reviewing or teaching ML fundamentals; selecting a learning paradigm or architecture; choosing hyperparameter-tuning tooling; designing an evaluation harness for either a classical or LLM workload; standing up production monitoring (drift, skew, retraining triggers); comparing 2026 frontier models for a build-vs-buy decision.
related_skills:
  - da-1-foundations-theory
  - da-1-3-probability-theory
  - da-1-4-statistical-inference-foundations
  - da-1-5-information-theory
  - da-1-6-epistemology-of-data
  - da-4-data-cleaning-preparation
  - da-6-statistical-modeling
  - da-8-data-visualization
  - da-9-reporting-communication
  - prompt-engineering
  - llm-context-engineering
  - rag-architecture
  - mongodb-atlas-vector-search
  - mongodb-search-ai
  - ai-datastores
---

# Machine Learning (Data Analysis Curriculum, Section 7)

Machine learning is the curriculum step where the analyst stops merely describing a sample and starts building a function that generalizes from data to unseen inputs. Section 6 (`da-6-statistical-modeling`) covered parametric statistical models grounded in explicit probabilistic assumptions. This section widens the lens to algorithms that learn flexible, often non-parametric mappings from data — and to the engineering scaffolding (tuning, evaluation, deployment, monitoring) that turns a trained model into a system that keeps working.

This skill is the curriculum reference for the seventh section of the data analysis path. It is intentionally broad: it sketches the territory and points to the deeper skills you should pull in for any specific build. The goal is to keep you from skipping the foundations — and from defaulting to the most fashionable architecture without a reason.

## Part 1 — Machine Learning Taxonomy

### 1.1 Three (now four) classical paradigms

ML divides historically into three paradigms; self-supervised learning has joined them as a first-class member since the rise of foundation models.

- **Supervised learning.** The training set consists of input–output pairs `(x, y)`. The model learns a function `f(x) ≈ y`, then predicts `y` for unseen `x`. Two sub-shapes:
  - *Classification* — `y` is discrete (spam vs. ham, malignant vs. benign, intent label).
  - *Regression* — `y` is continuous (house price, expected revenue, lead-time-to-failure).
- **Unsupervised learning.** No labels. The algorithm discovers structure: clusters (k-means, DBSCAN, HDBSCAN), latent topics (LDA, NMF), lower-dimensional manifolds (PCA, t-SNE, UMAP, autoencoders), density (KDE, Gaussian mixtures), or anomalies (isolation forest, one-class SVM).
- **Reinforcement learning (RL).** An agent interacts with an environment, observes a state, picks an action, receives a reward. Optimizes a policy `pi(a | s)` that maximizes expected discounted reward. Modern flavors include policy-gradient (PPO, GRPO), value-based (DQN, Rainbow), and model-based RL. RLHF and RLAIF are RL paradigms repurposed to align LLMs.
- **Self-supervised learning.** Labels are constructed automatically from the data itself (predict the next token, predict a masked patch, predict whether two augmentations of an image came from the same source). This is the engine behind every modern foundation model: BERT (masked language modeling), GPT/Claude/Gemini (autoregressive next-token prediction), DINO/MAE (masked image modeling), wav2vec (audio), CLIP (image-text contrastive).

Two other shapes appear at the edges:

- **Semi-supervised learning.** A small labeled set plus a much larger unlabeled set. Useful when labeling is expensive (medical imaging, legal review).
- **Active learning.** The model chooses which examples to send to a human labeler next — typically the examples it is most uncertain about — to maximize information gain per labeling dollar.

### 1.2 The bias-variance tradeoff

The bias-variance decomposition of expected prediction error on a fresh test point is:

`E[(y - f_hat(x))^2] = Bias[f_hat(x)]^2 + Var[f_hat(x)] + sigma^2`

- **Bias** — error from oversimplifying the true function. A linear model on a quadratic relationship has high bias and underfits.
- **Variance** — error from the model being too sensitive to the particular training sample. A 50-degree polynomial on 200 points has high variance and overfits.
- **Irreducible error (`sigma^2`)** — noise inherent in the data that no model can remove.

Increasing model capacity (more parameters, deeper trees, more layers) typically *lowers bias and raises variance.* Regularization, more data, ensembling, and architectural constraints all reduce variance at the cost of some bias.

Caveats:

- The classical U-shaped curve is an idealization. In *over-parameterized* regimes — the regime that virtually all modern deep learning lives in — test error often follows the "double descent" curve: it goes up as you cross the interpolation threshold, then *decreases again* as capacity grows further. Don't assume "more parameters = more overfitting" in 2026; that intuition came from a regime we have largely left behind.
- The decomposition is exact for squared-error loss. For 0/1 loss and cross-entropy the analog is messier; the qualitative tradeoff still holds.

### 1.3 Regularization toolbox

- **L1 (Lasso).** Adds `lambda * sum|w|` — drives weights to exactly zero, performs feature selection.
- **L2 (Ridge / weight decay).** Adds `lambda * sum(w^2)` — shrinks weights smoothly, prevents any single weight from dominating. The default in deep learning, usually applied to non-bias parameters only.
- **ElasticNet.** Convex combination of L1 and L2.
- **Dropout.** Stochastically zeros activations during training — equivalent to averaging an exponentially large ensemble at test time.
- **Early stopping.** Halt training when validation loss stops improving. Implicit regularization.
- **Data augmentation.** Cheap, often the single most effective regularizer when data is the bottleneck.
- **Label smoothing.** Replace one-hot targets with `(1 - epsilon) * onehot + epsilon / K` — reduces overconfidence.
- **Batch / Layer / RMS norm.** Normalization layers act as a side-channel regularizer in addition to their optimization role.

### 1.4 The classical model zoo

Worth knowing before reaching for deep learning:

- **Linear / logistic regression.** First-principles baseline. Cheap, interpretable, calibrates well, hard to beat on small tabular datasets.
- **Decision trees, random forests, gradient-boosted trees.** XGBoost, LightGBM, CatBoost still win most tabular ML competitions in 2026. Deep learning has not displaced GBTs on structured data.
- **Support vector machines.** Maximum-margin classifiers. Useful baseline on small-to-medium datasets, especially with the kernel trick.
- **k-Nearest Neighbors.** Lazy learner, no training phase, slow at inference. Useful as a sanity-check baseline and inside vector-search systems.
- **Naive Bayes.** Strong assumption of conditional independence; the assumption is almost always wrong; the model is almost always surprisingly good on text.

If your problem is on rows-and-columns tabular data, *start with a gradient-boosted tree*. If it's images, audio, or text, jump to Part 2.

## Part 2 — Deep Learning Fundamentals

Deep learning is the subset of ML that uses neural networks with multiple hidden layers. The defining shift versus classical ML is *representation learning*: the model learns its own features from raw data instead of consuming hand-engineered ones.

### 2.1 Convolutional Neural Networks (CNNs)

Specialized for grid-structured data (images, video frames, spectrograms).

- **Convolutional layer.** Slides a small learnable kernel over the input; weight sharing makes the layer translation-equivariant and dramatically reduces parameter count vs. dense layers.
- **Pooling.** Downsamples (max or average pool). Increases receptive field and adds approximate translation invariance.
- **Hierarchy.** Early layers learn edges and textures; middle layers learn parts; deep layers learn objects.
- **Canonical architectures.** LeNet → AlexNet → VGG → ResNet (residual connections, the unlock for >100 layers) → EfficientNet (compound scaling) → ConvNeXt (modernized CNN that competes with vision transformers).
- **Still relevant in 2026.** CNNs remain competitive on edge devices, medical imaging with limited data, and as efficient backbones inside hybrid CNN-Transformer models.

### 2.2 Recurrent Neural Networks (RNNs)

Built for sequences. State `h_t = f(x_t, h_{t-1})` is updated step-by-step.

- **Vanilla RNNs.** Suffer from vanishing and exploding gradients; struggle to capture long-range dependencies.
- **LSTM (Long Short-Term Memory).** Gated cell that learns when to remember and when to forget; the workhorse of 2014–2018.
- **GRU (Gated Recurrent Unit).** Simpler than LSTM, comparable performance, fewer parameters.
- **Bidirectional variants.** Process the sequence forward and backward; useful when the entire sequence is available at inference.

RNNs have been largely displaced by Transformers for text and most audio. They remain relevant for:

- Streaming inference where Transformer KV cache costs grow linearly with sequence length.
- Tiny edge models for time-series sensors.
- State-space models (Mamba, S4, S6) — the modern resurgence of recurrent computation, which scales linearly in sequence length and is closing the gap with Transformers on language modeling.

### 2.3 Why architecture choice still matters

The 2026 default for almost everything is "Transformer," but the right call depends on:

- **Inductive bias.** CNNs encode locality and translation equivariance; Transformers do not. On small datasets, the right inductive bias matters more than raw capacity.
- **Sequence length and inference cost.** Transformer attention is quadratic in sequence length. For >100k-token contexts, mixed approaches (sliding window, sparse attention, linear attention, Mamba-style SSMs) win on throughput.
- **Latency budget.** A 7B-parameter Transformer is overkill for keyword-spotting on a watch. A 1M-parameter CNN is the right call.

The next section covers Transformers in their own right — they earned it.

## Part 3 — Transformers and Foundation Models

### 3.1 The Transformer

Introduced in *Attention Is All You Need* (Vaswani et al., 2017), the Transformer replaced recurrence with self-attention. The mechanism:

- For each token, compute three vectors — query `Q`, key `K`, value `V` — via learned linear projections.
- Attention weights: `softmax(Q K^T / sqrt(d_k))` — measures how much each token should attend to every other token.
- Output: weighted sum of value vectors.
- *Multi-head* attention runs several attention operations in parallel on lower-dimensional projections, then concatenates.

Why this matters:

- **Parallelism.** Unlike RNNs, every position is processed simultaneously during training — drastically faster on GPUs and TPUs.
- **Long-range dependencies.** Any token can attend to any other in one layer; no vanishing gradients across the sequence axis.
- **Modality-agnostic.** The same block works for text (BERT, GPT), images (ViT), audio (Whisper), proteins (AlphaFold 2/3), and multi-modal inputs.

### 3.2 Foundation model paradigm

A foundation model is a large model pretrained on broad data with self-supervision, then *adapted* to downstream tasks via:

- **Zero-shot prompting** — describe the task in natural language.
- **Few-shot / in-context learning** — provide several examples in the prompt.
- **Fine-tuning** — gradient updates on task-specific data. Variants: full fine-tune, LoRA / QLoRA (low-rank adapters), DPO / KTO (preference tuning), RLHF, RLAIF.
- **Retrieval-augmented generation (RAG)** — combine the model with a vector store; see `rag-architecture`.
- **Tool use / function calling** — let the model invoke external APIs and process the results.

The 2018–2023 era was "pretrain a foundation model, fine-tune it for your task." The 2024–2026 era is increasingly "pretrain or buy a foundation model, then build *agent-shaped* systems around it that combine prompting, retrieval, tool use, and short structured fine-tuning."

## Part 4 — Frontier LLM Landscape (May 2026)

The 2026 frontier is multi-polar. No single model is best at everything; build pipelines that route to the right model for the request.

| Model | Provider | Released | Strengths | Notes |
| --- | --- | --- | --- | --- |
| Claude Opus 4.7 | Anthropic | Apr 2026 | Leading agentic coding, SWE-bench Verified 87.6%, 3.75MP image input | Strongest model for long-horizon coding agents; tuned for tool use |
| GPT-5.5 | OpenAI | Apr 2026 | Top Intelligence Index overall; first ground-up rebuild since GPT-4.5 | New pretraining corpus and objectives; broad capability across reasoning and writing |
| Gemini 3.1 Pro | Google DeepMind | 2026 | Leading scientific reasoning; multi-modal (image, audio, video) | Strong on long-context and grounded retrieval |
| Llama 4 Scout | Meta | 2026 | Open weights; 10M-token context (largest in any model) | Open-weight default for self-hosted; trades some quality for openness |
| Llama 4 / Muse Spark | Meta | Apr 2026 | Intelligence Index ~52 — within striking distance of frontier closed models | Open-weight contender |
| DeepSeek V3.2 | DeepSeek | 2026 | Best value-per-dollar at the frontier | High-volume budget workloads |
| Grok 4 | xAI | 2026 | Leads raw SWE-bench scores | Competitive on agentic coding |
| GLM-5.1 | Zhipu | 2026 | Briefly held #1 on SWE-bench Pro | Strongest open-source coding model at release |

Selection heuristics for May 2026:

- **Autonomous coding / agentic loops** → Claude Opus 4.7 (or 4.8 max when accessible).
- **General reasoning, writing, and multi-step problem solving** → GPT-5.5.
- **Multimodal (image + video + audio + text)** → Gemini 3.1 Pro.
- **Self-hosting, data sovereignty, or fine-tuning** → Llama 4 (Scout for long context, Maverick / Muse Spark for general use).
- **High-volume, latency-sensitive, cost-bound** → DeepSeek V3.2 or smaller specialist models (Haiku-class, Mini-class, Flash-class).

The model-switching pattern is now standard. Frameworks like LiteLLM, OpenRouter, and Anthropic / OpenAI's own router endpoints make per-request model selection cheap. Treat the LLM as a *replaceable component*, not a vendor commitment.

## Part 5 — Hyperparameter Tuning

Hyperparameters are the knobs *not* learned by gradient descent: learning rate, batch size, optimizer choice, weight decay, dropout rate, tree depth, number of estimators, kernel parameters. They dominate model quality — and they cost compute.

### 5.1 Search strategies

- **Manual / expert.** Reasonable starting point; never the final answer for anything you care about.
- **Grid search.** Exhaustive on a Cartesian product. Pathologically expensive in high dimensions. Only sensible when ≤3 hyperparameters with small discrete ranges.
- **Random search.** Sample uniformly from the search space. Bergstra & Bengio (2012) showed it strictly dominates grid search when only a few hyperparameters actually matter — which is most of the time.
- **Bayesian optimization (BO).** Fits a surrogate model (Gaussian process or TPE) over the objective and uses an acquisition function (Expected Improvement, UCB) to pick the next configuration. Strong default for 10–100 trials when each trial is expensive.
- **HyperBand.** Successive halving — start many trials with tiny budget, kill the worst, double the budget on survivors. Excellent for deep learning.
- **BOHB.** Bayesian Optimization + HyperBand. Bayesian early-stopping. Often the production-quality default.
- **PBT (Population-Based Training).** Train a population in parallel; periodically copy weights from top performers to bottom and perturb hyperparameters. Discovers schedules, not just static configs.

### 5.2 Tooling (2026)

- **Optuna.** Python-native, ergonomic define-by-run API, TPE sampler by default, integrated pruning. The default for individual researchers and most production ML teams.
- **Ray Tune.** Distributed-first, scales to large clusters, supports every algorithm above (PBT, BOHB, HyperBand, BO). The default when you have parallel compute.
- **Weights & Biases Sweeps.** Tightly integrated with experiment tracking; convenient when W&B is already your tracking layer.
- **KerasTuner, AutoGluon, FLAML.** Higher-level AutoML wrappers.
- **Vizier (Google) and SigOpt (Intel).** Hosted services.

### 5.3 Best-practice cheatsheet

- Always tune on a held-out validation set; never on the test set.
- Use *nested cross-validation* if your dataset is small enough that a single split is noisy.
- Log-uniform sample learning rate and weight decay; linear sample dropout and momentum.
- Tune learning rate first — it dominates everything else.
- Cap wall-clock budget per trial; let HyperBand or pruning kill stragglers.
- Re-tune when you change the dataset, the architecture family, or the optimizer.
- Save the full trial history; you will want to do a post-hoc importance analysis (Optuna provides `get_param_importances`).

## Part 6 — Evaluation Methodology

A model with no evaluation harness is a guess. The harness is the model.

### 6.1 Splits and resampling

- **Train / validation / test.** Train fits parameters; validation tunes hyperparameters and selects models; test is touched *once*, after every other decision is locked.
- **k-fold cross-validation.** Rotate through `k` train/val partitions; average the score. Standard for moderate-sized datasets.
- **Stratified k-fold.** Preserve class balance per fold. Mandatory for imbalanced classification.
- **Group / leave-one-group-out k-fold.** When rows have group identity (patient ID, customer ID), leakage between groups will inflate scores. Always group-split when groups exist.
- **Time-series splits.** Past predicts future. Use `TimeSeriesSplit` or rolling-origin evaluation; never shuffle a time series.
- **Nested CV.** Outer loop estimates generalization; inner loop tunes hyperparameters. Mandatory when the dataset is small and the model has many hyperparameters.

### 6.2 Classification metrics

For binary classification with a confusion matrix `(TP, FP, TN, FN)`:

- **Accuracy = (TP + TN) / N.** Misleading on imbalanced data — a "predict the majority" classifier hits 99% accuracy when 99% of cases are negative.
- **Precision = TP / (TP + FP).** Of predicted positives, how many were right. Optimize when false positives are expensive (spam filter sending real mail to junk).
- **Recall (sensitivity, TPR) = TP / (TP + FN).** Of actual positives, how many you caught. Optimize when false negatives are expensive (missing a cancer).
- **F1 = 2 * (P * R) / (P + R).** Harmonic mean. Balances precision and recall; useful default for imbalanced classes.
- **F-beta.** Weight recall over precision by factor beta (or vice versa). Pick beta from business cost.
- **ROC curve and AUC.** Plot TPR vs. FPR across thresholds; the area summarizes ranking quality independent of any single threshold. Insensitive to class balance.
- **Precision-Recall curve and PR-AUC.** Preferred over ROC when the positive class is rare.
- **Log-loss / cross-entropy.** Proper scoring rule; penalizes confidently-wrong predictions hard. The right metric when calibrated probabilities matter.
- **Calibration plots and Brier score.** Check whether predicted probabilities match observed frequencies. Calibrate with Platt scaling or isotonic regression if not.

For multi-class:

- **Macro F1** averages per-class F1 (treats all classes equally; punishes failure on rare classes).
- **Weighted F1** averages per-class F1 weighted by support (closer to accuracy).
- **Micro F1** equals accuracy in single-label multi-class. Different in multi-label.

### 6.3 Regression metrics

- **MAE (Mean Absolute Error).** Robust to outliers; same units as the target.
- **RMSE (Root Mean Squared Error).** Same units as the target; penalizes large errors more than MAE.
- **MAPE (Mean Absolute Percentage Error).** Scale-free; undefined when actuals are zero.
- **sMAPE.** Symmetric MAPE — better-behaved near zero.
- **R^2.** Variance explained relative to a mean baseline. Can be negative if your model is worse than the mean.
- **Quantile loss / pinball loss.** When you need quantile forecasts (probabilistic forecasting).

### 6.4 Ranking and recommender metrics

- **MAP@k, NDCG@k, MRR.** Use when output order matters (search, recommendations).
- **Hit-rate@k.** Did the true item appear in the top-k.

### 6.5 LLM evaluation

LLMs broke classical evaluation: outputs are free-form text, ground truth is fuzzy, and the most interesting capabilities (reasoning, multi-step planning, tool use) don't fit a confusion matrix.

- **MMLU (Massive Multitask Language Understanding).** 14,042 four-option multiple choice questions across 57 subjects. The historical default; now saturated — frontier models score >90% and the metric stops discriminating. Use only as a sanity check.
- **MT-Bench.** Multi-turn benchmark with harder, reasoning-heavy questions. Uses LLM-as-judge (typically GPT-4 class or higher) for scoring.
- **HELM (Holistic Evaluation of Language Models).** Stanford CRFM framework. Scores across accuracy, calibration, robustness, fairness, bias, toxicity, and efficiency. The right framework for risk-aware evaluation.
- **Chatbot Arena.** Crowd-sourced pairwise preference voting; Elo ratings. The closest thing to a ground-truth signal of which model people actually prefer.
- **SWE-bench / SWE-bench Verified / SWE-bench Pro.** Real GitHub issues with hidden test suites. The frontier benchmark for agentic coding; Claude Opus 4.7 currently leads SWE-bench Verified at 87.6%.
- **HumanEval, MBPP, LiveCodeBench.** Programming benchmarks of varying difficulty and contamination risk.
- **GPQA, ARC-AGI, FrontierMath.** "Hard" reasoning benchmarks that have not saturated.
- **Domain-specific evals.** MedQA, LegalBench, FinanceBench, MMMU for multimodal. Build your own when your domain is the constraint.

### 6.6 LLM-as-judge

Free-form outputs are commonly evaluated by a stronger LLM scoring against a rubric.

- **Pros.** Fast, scales, captures qualitative judgments humans agree with.
- **Cons.** Position bias (judges prefer the first answer); self-preference (a model judging itself scores itself higher); verbosity bias (judges prefer longer answers); rubric drift.
- **Mitigations.** Randomize position; use a different judge family than the model under test; calibrate judge against a small human-labeled gold set; use pairwise comparison rather than absolute scoring when possible.
- **Industry trend (2026).** Multi-judge consensus and structured rubrics (G-Eval, DeepEval, Patronus, Braintrust) have become standard for production LLM evaluation.

## Part 7 — MLOps

The seam between "model that worked once in a notebook" and "model that keeps working in production."

### 7.1 Experiment tracking

- **MLflow.** Open-source, OSS-friendly default. Tracks parameters, metrics, artifacts, models. Self-hosted or Databricks-hosted.
- **Weights & Biases.** Commercial; richer visualization, better collaboration; standard at well-funded teams.
- **Neptune.ai, Comet.** Commercial alternatives with different tradeoffs.
- **DVC + Git.** When the team is small and the model is small, version data and code together in Git/DVC; no separate tracking layer needed.

Whatever you pick, log: code commit hash, dataset hash, hyperparameters, environment (Python, library versions), training metrics, validation metrics, the trained artifact, and the evaluation report.

### 7.2 Data and concept drift

- **Covariate shift.** Input distribution `P(x)` changes; the relationship `P(y | x)` is stable.
- **Concept drift.** The relationship `P(y | x)` changes — customer preferences evolve, fraud patterns adapt, language drifts.
- **Label drift.** `P(y)` changes — class balance shifts (more fraud after a policy change).

Detection methods:

- **Population Stability Index (PSI), Jensen-Shannon divergence, KL divergence.** Per-feature distribution comparisons.
- **Kolmogorov-Smirnov, chi-squared, Wasserstein.** Statistical tests for distribution shift.
- **Performance monitoring.** When labels arrive late or never, monitor proxies (prediction distribution, model confidence) instead of accuracy directly.
- **Tooling.** Evidently AI, WhyLabs, Arize, Fiddler. Evidently is the strongest OSS choice.

Industry data point (Weights & Biases, 2023): 62% of organizations experience meaningful model degradation within 12 months when not actively monitored.

### 7.3 Train-serve skew

Different from drift: drift is the world changing, skew is *your* code being inconsistent between training and serving.

- **Feature parity.** Train and serve must use byte-identical preprocessing. The most common production bug. Use a feature store (Feast, Tecton, Hopsworks, or platform-native) so both code paths consume the same transformation logic.
- **Schema parity.** Column types, null handling, encoding (categorical hash → integer mapping). Pin them with a schema check at both training and serving boundaries (TFX SchemaGen, Great Expectations, Pandera).
- **Lookup parity.** If you compute `customer_lifetime_value` from raw events at training time, you cannot read a stale rolled-up cache at serving time. Either materialize the same way in both, or push the join logic into the feature store.
- **Time-leakage.** A feature available at training time but not at prediction time (e.g., a label-correlated event happened after the prediction would have been made). Build features with explicit "as-of" timestamps.

### 7.4 Deployment patterns

- **Shadow deployment.** New model runs in parallel; predictions logged, not served.
- **Canary deployment.** Send X% of traffic to the new model; ramp up if metrics hold.
- **A/B test.** Randomized assignment with statistical analysis; tied to a business metric, not just model accuracy.
- **Multi-armed bandit.** Adaptive allocation; converges to the best variant faster than fixed A/B.
- **Champion / challenger.** Production model is the champion; new models challenge on a holdout stream; promote when challenger beats champion by a stat-sig margin.

### 7.5 Retraining triggers

- **Scheduled.** Daily, weekly, monthly. Cheap to operate; misses fast-moving drift.
- **Drift-triggered.** Retrain when drift detectors fire. Avoids retraining when nothing changed.
- **Performance-triggered.** Retrain when monitored metrics fall below threshold. Requires fast-arriving labels.
- **Continuous training.** Stream-based; the model updates on every batch. Standard at top-tier ads and recommendation companies.

### 7.6 Reproducibility checklist

- Pinned dependency versions (`requirements.lock`, `poetry.lock`, `uv.lock`, `conda-lock`).
- Hashed dataset snapshots (DVC, LakeFS, Delta Lake).
- Fixed random seeds — but understand they don't guarantee bit-identical runs on GPUs.
- Containerized training (Docker, Singularity / Apptainer).
- A recorded `mlflow run` or equivalent with everything above attached.

## Anti-Patterns to Avoid

- **Tuning hyperparameters on the test set.** You no longer have a test set; you have a second validation set. Every leak inflates the reported number.
- **Reporting only mean accuracy on an imbalanced dataset.** Report per-class metrics, the confusion matrix, and a relevant aggregate (macro F1, PR-AUC).
- **Comparing models on different splits.** Always cross-validate over the same splits.
- **Shuffling a time series before splitting.** Catastrophic leakage. Always time-split temporal data.
- **Ignoring calibration in classification.** A model that says "90% confident" should be right 90% of the time. If it isn't, downstream decisions built on those probabilities are broken.
- **Picking the model with the highest validation score and shipping it.** Validate the pipeline (data + features + model + post-processing), not the model in isolation.
- **Building a "vibes-only" LLM eval.** "Looks better to me" is not a metric. Build a holdout set, a rubric, and a judge — even a small one.
- **Treating the LLM as a fixed dependency.** The frontier moves every 2–3 months. Build for model-substitutability from day one.
- **Skipping monitoring because "it works in dev."** Drift is the default state of the world. A model without a drift monitor is a model that will silently fail.

## Related Skills

This curriculum entry intentionally cross-links into adjacent skills rather than duplicating their depth. Recommended next reads:

- `da-1-foundations-theory` — measurement theory and probability that anchor everything in this skill.
- `da-1-3-probability-theory` and `da-1-4-statistical-inference-foundations` — the probabilistic substrate behind every loss function and confidence interval here.
- `da-1-5-information-theory` — KL divergence, entropy, cross-entropy. Useful when drift detection or model loss functions are the focus.
- `da-1-6-epistemology-of-data` — correlation vs causation, reproducibility, replicability. Critical before claiming a model "found" anything.
- `da-4-data-cleaning-preparation` — feature engineering, encoding, imputation. These choices dominate model quality more than algorithm choice on most tabular problems.
- `da-6-statistical-modeling` — the parametric-model sibling of this skill. Read alongside.
- `da-8-data-visualization` — diagnostics, residual plots, ROC curves, learning curves.
- `da-9-reporting-communication` — communicating model results to non-technical stakeholders.
- `prompt-engineering`, `llm-context-engineering`, `phe`, `ph`, `prompt-deep-optimizer` — when the "model" is an LLM and your knob is the prompt.
- `rag-architecture` — when foundation-model output must be grounded in your own corpus.
- `mongodb-atlas-vector-search`, `mongodb-search-ai`, `mongodb-atlas-search` — building the retrieval substrate for RAG and semantic search on Atlas.
- `ai-datastores`, `ai-languages`, `llm-models` — vendor- and tool-side reference material.
- `mongodb-atlas-stream-processing` — when training or scoring runs in stream rather than batch.

## References

- Bergstra, J. and Bengio, Y. (2012). *Random Search for Hyperparameter Optimization.* JMLR.
- Vaswani, A. et al. (2017). *Attention Is All You Need.* NeurIPS.
- Liang, P. et al. (2022). *Holistic Evaluation of Language Models (HELM).* Stanford CRFM.
- Chiang, W.-L. et al. (2023). *Chatbot Arena: An Open Platform for Evaluating LLMs by Human Preference.*
- Belkin, M. et al. (2019). *Reconciling Modern Machine Learning Practice and the Bias-Variance Trade-off.* PNAS (double descent).
- Hendrycks, D. et al. (2021). *Measuring Massive Multitask Language Understanding (MMLU).* ICLR.
- Jimenez, C. et al. (2024). *SWE-bench: Can Language Models Resolve Real-World GitHub Issues?* ICLR.
- LM Council, Vellum, Artificial Analysis — live LLM leaderboards (May 2026).
- Anthropic, OpenAI, Google DeepMind, Meta, DeepSeek, xAI, Zhipu — model release notes, April–May 2026.
- Evidently AI, Weights & Biases — MLOps drift and monitoring industry reports (2023–2026).


