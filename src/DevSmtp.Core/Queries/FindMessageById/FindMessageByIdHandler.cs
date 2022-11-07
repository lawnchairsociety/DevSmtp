using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Queries
{
    public sealed class FindMessageByIdHandler : IQueryHandler<FindMessageById, FindMessageByIdResult>
    {
        private readonly IDataStore _dataStore;

        public FindMessageByIdHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public async Task<FindMessageByIdResult> ExecuteAsync(FindMessageById query, CancellationToken cancellationToken = default)
        {
            cancellationToken.ThrowIfCancellationRequested();

            try
            {
                var message = await this._dataStore.FindByIdAsync(query.Id, cancellationToken);
                return new FindMessageByIdResult(message);
            }
            catch (OperationCanceledException)
            {
                throw;
            }
            catch (Exception ex)
            {
                var message = $"Failed to find message by '{query.Id.Value}'.";
                var error = new FindMessageByIdException(message, ex);

                return new FindMessageByIdResult(error);
            }
        }
    }
}
