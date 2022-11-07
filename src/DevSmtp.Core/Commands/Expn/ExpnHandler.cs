using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class ExpnHandler : ICommandHandler<Expn, ExpnResult>
    {
        private readonly IDataStore _dataStore;

        public ExpnHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<ExpnResult> ExecuteAsync(Expn command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
