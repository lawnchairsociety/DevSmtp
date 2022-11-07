using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Queries
{
    public sealed class GetMessagesHandler : IQueryHandler<GetMessages, GetMessagesResult>
    {
        private readonly IDataStore _dataStore;

        public GetMessagesHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public async Task<GetMessagesResult> ExecuteAsync(GetMessages query, CancellationToken cancellationToken = default)
        {
            cancellationToken.ThrowIfCancellationRequested();

            try
            {
                var results = await this._dataStore.GetAsync(cancellationToken);
                return new GetMessagesResult(results);
            }
            catch (OperationCanceledException)
            {
                throw;
            }
            catch (Exception ex)
            {
                var message = "Failed to get messages.";
                var error = new GetMessagesException(message, ex);

                return new GetMessagesResult(error);
            }
        }
    }
}
